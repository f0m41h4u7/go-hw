package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Prev  *listItem
	Next  *listItem
}

type list struct {
	front  *listItem
	back   *listItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushBack(v interface{}) *listItem {
	newItem := &listItem{
		Value: v,
		Prev:  nil,
		Next:  l.back,
	}

	if l.back == nil {
		l.front = newItem
	} else {
		l.back.Prev = newItem
	}

	l.back = newItem
	l.length++
	return newItem
}

func (l *list) PushFront(v interface{}) *listItem {
	newItem := &listItem{
		Value: v,
		Prev:  l.front,
		Next:  nil,
	}

	if l.front == nil {
		l.back = newItem
	} else {
		l.front.Next = newItem
	}
	l.front = newItem
	l.length++
	return newItem
}

func (l *list) Remove(i *listItem) {
	if i == nil {
		return
	}
	l.length--

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.back = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.front = i.Prev
	}
}

func (l *list) MoveToFront(i *listItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{}
}
