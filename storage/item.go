package storage

const NoExpire = -1

// Storage is core element, which consist data in map[string]interface{}
// format.
type ItemInterface interface {
	Value() interface{}
	SetValue(interface{})
	Key() string
	Expire() (int, bool)
	SetExpire(i int)
	String() string
}

type Item struct {
	ItemInterface
	key string
	value interface{}
	expire int
}

func NewItem(key string, val interface{}) ItemInterface {
	return &Item{key: key, value: val}
}

func (n *Item) Key() string { return n.key }
func (n *Item) Value() interface{} { return n.value }
func (n *Item) SetValue(v interface{}) { n.value = v }

func (n *Item) SetExpire(e int) {
	if e < 1 {
		n.expire = NoExpire
	}
	n.expire = e
}

func (n *Item) Expire() (int, bool) {
	if n.expire == NoExpire {
		return NoExpire, false
	}
	return n.expire, true
}

