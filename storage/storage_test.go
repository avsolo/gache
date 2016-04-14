package storage

import (
	"fmt"
	"time"
	"testing"
	"strings"
	"math/rand"
	"crypto/md5"
	"github.com/stretchr/testify/assert"

	"avsolo/gache/lib"
	"avsolo/gache/storage"
)

var log = lib.NewLogger("")
var _s = fmt.Sprintf
var yota int = 1

type TestItem struct {
	Field interface{}
}

func NewTestItem(val interface{}) *TestItem {
	return &TestItem{Field: val}
}

// makeTestItem generate and return new Key, Value, TTL
func makeTestItem() (key string, item *TestItem, ttl int) {
	rand.Seed(time.Now().Unix())
	ttl = rand.Intn(10) + 1
	hash := md5.Sum([]byte(fmt.Sprintf("%d", ttl)))
	yota++
	key = fmt.Sprintf("%x%d", hash, yota)
	item = NewTestItem(strings.Repeat(key, 10))
	return
}

// Test simple get/set for struct{}
func TestSetSuccess(t *testing.T) {
	// Ok item
	s := storage.NewStorage()
	keyOk, itemOk, ttl := makeTestItem()
	// Wrong item
	_, itemWrong, ttlWrong := makeTestItem()

	// Set
	err := s.Set(keyOk, itemOk, ttl)
	assert.Nil(t, err, _s("Set error: %v", err))

	// Check that we can't set againt same key (only update)
	err = s.Set(keyOk, itemWrong, ttlWrong)
	assert.Equal(t, err, storage.ErrAlreadyExists, _s("Second set error: %v", err))

	// Check that we can get our TestItem back from Storage
	itemGet, err := s.Get(keyOk)
	assert.Nil(t, err)
	assert.Equal(t, itemOk, itemGet, "Value was changed by second Set")

	// Check gotten item with our Test item field by field
	itemGetConv, _ := itemGet.(*TestItem)
	assert.Equal(t, itemOk.Field, itemGetConv.Field, "Cant' convert gotten interface to TestItem")
}

func TestUpdateSuccess(t *testing.T) {
	// Prepare
	s := storage.NewStorage()
	key, item, ttl := makeTestItem()
	_ = s.Set(key, item, ttl)

	// Update same key
	_, itemUpd, ttlUpd := makeTestItem()
	err := s.Update(key, itemUpd, ttlUpd)
	assert.Nil(t, err, "Update error: %v", err)

	// Check after update
	itemGet, err := s.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, itemUpd, itemGet, "Value not updated")
}

func TestDeleteSuccess(t *testing.T) {
	s := storage.NewStorage()
	key, item, ttl := makeTestItem()
	_ = s.Set(key, item, ttl)

	// Delete same key
	s.Delete(key)

	// Check after delete
	itemGet, err := s.Get(key)
	assert.Equal(t, storage.ErrNotFound, err)
	assert.Nil(t, itemGet, "Value not deleted")
}

// Expire /////////////////////////////////////////////////////////////////////

func TestExpireSuccess(t *testing.T) {
	s := storage.NewStorage()
	key, item, _ := makeTestItem()
	ttl := 2

	// Set, and check while key expired
	_ = s.Set(key, item, ttl)

	// Wait and check that our key not exists
	time.Sleep(time.Duration(ttl + 1) * time.Second)
	itemGet, err := s.Get(key)
	assert.Equal(t, storage.ErrNotFound, err, "Wrong error type while returning expired item")
	assert.Nil(t, itemGet, "Value was changed (must to be expired)")
}

// Lists //////////////////////////////////////////////////////////////////////

func TestSetTTLSuccess(t *testing.T) {
	s := storage.NewStorage()
	key, item, _ := makeTestItem()
	ttl := 2

	// Set first TTL
	_ = s.Set(key, item, ttl)

	// Set new TTL
	newTtl := 4
	newTtlStamp := storage.MakeTTLStamp(newTtl)
	err := s.SetTTL(key, newTtl)
	assert.Nil(t, err, "SetTTL failed")

	// Check that TTL updated
	updatedTtl, err := s.GetExpire(key)
	assert.Equal(t, newTtlStamp, updatedTtl, "Updated TTL not equals its value")
}

func TestLSetSuccess(t *testing.T) {
	// Prepare
	s := storage.NewStorage()
	keyOk, _, ttl := makeTestItem()
	valOk := []string{ "value1", "value2", "value3" }

	// Set
	err := s.LSet(keyOk, valOk[0], valOk[1], valOk[2], ttl)
	assert.Nil(t, err, _s("LSet error: %v", err))

	// We can't set again
	err = s.LSet(keyOk, valOk[2], valOk[1], valOk[0], ttl)
	assert.Equal(t, storage.ErrAlreadyExists, err)

	// Check after update
	getList, err := s.LGet(keyOk)
	assert.Nil(t, err)
	convertedList, ok := getList.(storage.ItemListInterface)
	assert.Equal(t, convertedList.Len(), len(valOk))
	assert.True(t, ok)

	// Check all elements
	for i := len(valOk)-1; i >= 0; i-- {
		el, ok := convertedList.Pop()
		assert.True(t, ok)
		assert.Equal(t, valOk[i], el)
	}
}

func TestLPushSuccess(t *testing.T) {
	// Prepare
	s := storage.NewStorage()
	keyOk, _, ttl := makeTestItem()
	valOk := []int{10}

	// Set
	err := s.LSet(keyOk, valOk[0], ttl)
	getList, err := s.LGet(keyOk)
	convertedList, _ := getList.(storage.ItemListInterface)
	assert.Equal(t, convertedList.Len(), len(valOk), "Lenth of firstly setted list invalid")

	// Add one to list
	newVal := 123
	err = s.LPush(keyOk, newVal)
	assert.Nil(t, err, "Push error")

	// Check new size
	getList, err = s.LGet(keyOk)
	convertedList, _ = getList.(storage.ItemListInterface)
	assert.Equal(t, convertedList.Len(), len(valOk) + 1, "Lenth after push new element invalid")

	// Pop just added
	varGet, err := s.LPop(keyOk)
	assert.Nil(t, err, "Error while Pop from list")
	assert.Equal(t, varGet, newVal, "LPop return wrong element (must be new pushed)")

	// Pop last
	varGet, err = s.LPop(keyOk)
	assert.Nil(t, err, "Error LPop element")
	assert.Equal(t, varGet, valOk[0], "LPop return wrong element (must be first added)")

	// No value left
	varGet, err = s.LPop(keyOk)
	assert.Equal(t, err, storage.ErrNotFound, "Empty list must return nil and error")
	assert.Nil(t, varGet, "List must be empty")
}

// Dicts //////////////////////////////////////////////////////////////////////

func TestDSetSuccess(t *testing.T) {
	// Prepare
	s := storage.NewStorage()
	rkey, _, ttl := makeTestItem()
	k1, k2, k3 := "key1",	"key2",   "key3"
	v1, v2, v3 := "value1", "value2", "value3"

	// Set
	err := s.DSet(rkey, k1,v1, k2,v2, k3,v3, ttl)
	assert.Nil(t, err, _s("DSet error: %v", err))

	// Check that all key exists in map
	checkMapKeys(t, s, rkey, map[string]interface{}{k1:v1, k2:v2, k3:v3})

	// We can't set again
	err = s.DSet(rkey, k2,v2, k1,v1, k1,v1, ttl)
	assert.Equal(t, storage.ErrAlreadyExists, err)
}

func TestDAddSuccess(t *testing.T) {
	s := storage.NewStorage()
	rkey, _, ttl := makeTestItem()
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := 10, 23, 123

	err := s.DSet(rkey, map[string]interface{}{k1:v1}, ttl)

	// Add new key
	err = s.DAdd(rkey, k2, v2)
	assert.Nil(t, err)

	err = s.DAdd(rkey, k3, v3)
	assert.Nil(t, err)

	// Check that added key exists
	checkMapKeys(t, s, rkey, map[string]interface{}{k3:v3, k2:v2, k1:v1})
}

func TestDDelSuccess(t *testing.T) {
	s := storage.NewStorage()
	rkey, _, ttl := makeTestItem()
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := 10, 23, 123

	err := s.DSet(rkey, k1,v1, k2,v2, k3,v3, ttl)

	// Delete key
	s.DDel(rkey, k2)

	// Check
	r, err := s.DGet(rkey, k2)
	assert.Equal(t, storage.ErrNotFound, err, "Found deleted value")
	assert.Nil(t, r)

	// Check that rest of keys exists
	checkMapKeys(t, s, rkey, map[string]interface{}{k3:v3, k1:v1})
}

func checkMapKeys(t *testing.T, s *storage.Storage, rkey string, data map[string]interface{}) {
	for k, v := range data {
		r, err := s.DGet(rkey, k)
		assert.Nil(t, err, _s("DGet error for key %s", k))
		assert.Equal(t, v, r, _s("Value for key %s incorrect", k))
	}
}
