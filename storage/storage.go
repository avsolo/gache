// Package storage provides core functionality of geep caching mechanism
// and consist API of exported methods
package storage

import (
	"sync"
	"time"
)

// Storage is core element, which consist data in map[string]interface{}
// format.
type Storage struct {
    lock sync.RWMutex
    data map[string]ItemInterface
	expire map[int]map[string]struct{}
}

// NewStorage create a new instance of Storage. You can create any number of
// Storage and all of them will be work separately. Also NewStorate start
// expire traking - 1 sec timer. See startTiker for more details.
func NewStorage() *Storage {
	s := &Storage{
		data: map[string]ItemInterface{},
		expire: map[int]map[string]struct{}{},
	}
	go s.startTicker()
	return s
}

///////////////////////////////////////////////////////////////////////////
// Base getters/setters for struct{}
///////////////////////////////////////////////////////////////////////////

// Set find key with same name and if not exist - create on. If key
// exists do nothing and return ErrAlreadyExists error
func (s *Storage) Set(key string, val interface{}, ttl int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, found := s.data[key]; found {
		return ErrAlreadyExists
	}
	return s.setUnsafe(key, val, ttl)
}

// setUnsafe sets key/value to data and TTL
func (s *Storage) setUnsafe(key string, val interface{}, ttl int) error {
	s.data[key] = NewItem(key, val)
	s.data[key].SetExpire(MakeTTLStamp(ttl))
	s.setTTLUnsafe(key, ttl)
	return nil
}

// Get finds and return key from Storage. It uses internal Go mechanism
// and return bool as second argument
func (s *Storage) Get(key string) (interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if d, found := s.data[key]; found {
		return d.Value(), nil
	}
	return nil, ErrNotFound
}

// Update finds key in Storage and set new val and ttl if key found
// If not, return ErrNotFound error
func (s *Storage) Update(key string, val interface{}, ttl int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, found := s.data[key]; !found {
		return ErrNotFound
	}
	return s.setUnsafe(key, val, ttl)
}

// Delete finds and delete key. Uses Go internal mechanism
func (s *Storage) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	el, found := s.data[key]
	if !found {
		return
	}
	if exp, ok := el.Expire(); ok {
		delete(s.expire[exp], el.Key())
	}
	delete(s.data, key)
}

///////////////////////////////////////////////////////////////////////////
// List getters/setters
///////////////////////////////////////////////////////////////////////////

// LSet validate and convert string to list before call Set
// last argument in args must be TTL (interer)
func (s *Storage) LSet(key string, args ...interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, found := s.data[key]; found {
		return ErrAlreadyExists
	}
	ttlRaw, args := args[len(args)-1], args[:len(args)-1]
	ttl, ok := ttlRaw.(int)
	if !ok {
		return ErrBadTTL
	}
	l := NewItemList()
	for _, val := range args {
		l.Push(val)
	}
	return s.setUnsafe(key, l, ttl)
}

// LGet return ItemListInterface which implements simle stack interface
func (s *Storage) LGet(key string) (ItemListInterface, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	d, found := s.data[key]
	if !found {
		return nil, ErrNotFound
	}
	if itemList, ok := (d.Value()).(ItemListInterface); ok {
		return itemList, nil
	}
	return nil, ErrNotList
}

// getListUnsafe returns data withot sync.RLock. This method must be used
// with care and in the same gorutine
func (s *Storage) getListUnsafe(key string) (ItemListInterface, error) {
	el, found := s.data[key]
	if !found {
		return nil, ErrNotFound
	}
	list, ok := (el.Value()).(ItemListInterface)
	if !ok {
		return nil, ErrNotList
	}
	return list, nil
}

// LPush validate and convert string to list before call Set
func (s *Storage) LPush(key string, val interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	el, found := s.data[key]
	if !found {
		log.Warnf("Key '%s' not found", key)
		return ErrNotFound
	}
	list, ok := (el.Value()).(ItemListInterface)
	if !ok {
		log.Warnf("Key '%s' not list", key)
		return ErrNotList
	}
	list.Push(val)
	return nil
}

// LPop validate and convert string to list before call Set
func (s *Storage) LPop(key string) (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	l, err := s.getListUnsafe(key)
	if err != nil {
		return nil, err
	}

	if res, found := l.Pop(); found {
		return res, nil
	}
	return nil, ErrNotFound
}

///////////////////////////////////////////////////////////////////////////////
// Maps
///////////////////////////////////////////////////////////////////////////////

// DSet is set map[string] hash as Storage key
func (s *Storage) DSet(key string, args ...interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, found := s.data[key]; found {
		return ErrAlreadyExists
	}
	ttlRaw, args := args[len(args)-1], args[:len(args)-1]
	ttl, ok := ttlRaw.(int)
	if !ok {
		return ErrBadTTL
	}
	argsLen := len(args)

	// If it's map[string] already
	if argsLen == 1 {
		oneMap, correct := args[0].(map[string]interface{})
		if !correct {
			return ErrBadMap
		}
		return s.setUnsafe(key, oneMap, ttl)
	}

	if argsLen % 2 != 0 {
		return ErrBadMap
	}
	initMap := map[string]interface{}{}
	for i := 0; i < argsLen; i += 2 {
		k, ok := args[i].(string)
		if !ok {
			return ErrBadMap
		}
		initMap[k] = args[i+1]
	}
	return s.setUnsafe(key, initMap, ttl)
}

// DGet return value by key/subkey from map
func (s *Storage) DGet(rkey, skey string) (interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	d, found := s.data[rkey]
	if !found {
		log.Warnf("RKey %s not found for dict", rkey)
		return nil, ErrNotFound
	}
	itemMap, ok := (d.Value()).(map[string]interface{})
	if ok {
		if val, found := itemMap[skey]; found {
			return val, nil
		}
		log.Warnf("SKey %s not found for dict", skey)
		return nil, ErrNotFound
	}
	log.Warnf("Key %s not dict", rkey)
	return nil, ErrNotDict
}

// DAdd adds value by key/subkey
func (s *Storage) DAdd(rkey, skey string, val interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	el, found := s.data[rkey]
	if !found {
		return ErrNotFound
	}
	hash, ok := el.Value().(map[string]interface{})
    if !ok {
        return ErrNotDict
	}
    hash[skey] = val
    el.SetValue(hash)
    s.data[rkey] = el
    return nil
}

// DDel remote value by key/subkey
func (s *Storage) DDel(rkey, skey string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	m, found := s.data[rkey]
	if !found {
		return
	}
	if itemMap, ok := (m.Value()).(map[string]interface{}); ok {
		if _, found := itemMap[skey]; found {
			delete(itemMap, skey)
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
// TTL and tiker
///////////////////////////////////////////////////////////////////////////////

// SetTTL is find and remove old expire value and set new one
func (s *Storage) SetTTL(key string, ttl int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.setTTLUnsafe(key, ttl)
}

func (s *Storage) DeleteTTL(exp int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.expire, exp)
	return nil
}

// setTTLUnsafe isn't set any thread lock while it's set expire value
func (s *Storage) setTTLUnsafe(key string, ttl int) error {
	el, found := s.data[key]
	if !found {
		return ErrNotFound
	}
	cur, err := s.getExpireUnsafe(key)
	if err == nil {
		delete(s.expire[cur], key) // Remove current expire value
	}
	el.SetExpire(MakeTTLStamp(ttl)) // Update expire in Item

	t := MakeTTLStamp(ttl)
	_, found = s.expire[t]
	if ! found {
		s.expire[t] = map[string]struct{}{}
	}
	s.expire[t][key] = struct{}{}
	return nil
}

// GetTTL returns amount of seconds till key will expire
func (s *Storage) GetExpire(key string) (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.getExpireUnsafe(key)
}

func (s *Storage) getExpireUnsafe(key string) (int, error) {
	el, found := s.data[key]
	if !found {
		return NoExpire, ErrNotFound
	}
	exp, ok := el.Expire()
	if !ok {
		return NoExpire, ErrNoExpire
	}
	return exp, nil
}

// startTicker is timer with 1 sec tick, which finds timestamp key in 
// ExpireList.expire and then delete keys from ExpireList.data which has expire time
func (s *Storage) startTicker() {
	ticker := time.NewTicker(1 * time.Second)
	for tick := range ticker.C {
		t := int(tick.Unix())
		items, found := s.expire[t]
		if !found {
			continue
		}
		for key, _ :=range items {
			delete(items, key)
			s.Delete(key)
		}
	}
}

// Flush recursively delete keys and exprire data from Storage
func (s *Storage) Flush() {
	for k, v := range s.data {
		delete(s.data, k)
		if exp, ok := v.Expire(); ok {
			delete(s.expire, exp)
		}
	}
}

// MakeTTLStamp calculate absolute timestamp from ttl seconds
func MakeTTLStamp(ttl int) int {
	return int(time.Now().Add(time.Duration(int64(ttl)) * time.Second).Unix())
}
