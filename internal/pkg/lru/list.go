package lru

// List is an interface that represents a linked list. It provides methods to manipulate the list and access its elements.
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// ListItem is a structure that represents an element in a linked list.
// It contains a value and pointers to the next and previous elements in the list.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

// NewList creates a new instance of the List interface.
func NewList() List {
	return new(list)
}

// Len returns the number of elements in the list.
func (l *list) Len() int {
	return l.len
}

// Front returns the first element in the list.
func (l *list) Front() *ListItem {
	return l.front
}

// Back returns the last element in the list.
func (l *list) Back() *ListItem {
	return l.back
}

// PushAfter adds a new item after the specified previous item in the list.
func (l *list) PushAfter(prev *ListItem, item *ListItem) {
	item.Prev = prev

	if prev.Next == nil {
		item.Next = nil
		l.back = item
	} else {
		item.Next = prev.Next
		prev.Next.Prev = item
	}

	prev.Next = item
	l.len++
}

// PushBefore adds a new item before the specified next item in the list.
func (l *list) PushBefore(next *ListItem, item *ListItem) {
	item.Next = next

	if next.Prev == nil {
		item.Prev = nil
		l.front = item
	} else {
		item.Prev = next.Prev
		next.Prev.Next = item
	}

	next.Prev = item
	l.len++
}

// PushItemFront adds a new item to the front of the list.
func (l *list) PushItemFront(i *ListItem) {
	if l.front == nil {
		l.front = i
		l.back = i
		i.Prev, i.Next = nil, nil
		l.len++

		return
	}

	l.PushBefore(l.front, i)
}

// PushFront creates a new item with specified value and adds it to the front of the list.
func (l *list) PushFront(v any) *ListItem {
	item := &ListItem{Value: v}
	l.PushItemFront(item)

	return item
}

// PushBack adds a new item with specified value and adds it to the back of the list.
func (l *list) PushBack(v any) *ListItem {
	item := &ListItem{Value: v}

	if l.back == nil {
		l.PushItemFront(item)
	} else {
		l.PushAfter(l.back, item)
	}

	return item
}

// Remove removes an item from the list.
func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

// MoveToFront moves an existing item to the front of the list.
func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	l.Remove(i)
	l.PushItemFront(i)
}
