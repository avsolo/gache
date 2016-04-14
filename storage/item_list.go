package storage

import "container/list"

type ItemList struct {
	ItemInterface
	items *list.List
}

type ItemListInterface interface {
	Push(v interface{}) *list.Element
	Pop() (interface{}, bool)
	Len() int
}

func NewItemList() *ItemList {
	return &ItemList{ items: list.New(), }
}

func (l *ItemList) Len() int {
	return l.items.Len()
}

func (l *ItemList) Push(v interface{}) *list.Element {
	return l.items.PushBack(v)
}

func (l *ItemList) Pop() (interface{}, bool) {
	if l.items.Len() == 0 {
		return nil, false
	}
	el := l.items.Back()
	return l.items.Remove(el), true
}

func (l *ItemList) Delete(v interface{}) *list.Element {
	return l.items.PushBack(v)
}
