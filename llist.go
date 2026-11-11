package list

import "sync"

// LElement is an element of a LList. All its methods are safe for concurrent
// use. Each LList method returns a new LElement wrapper, so compare elements
// by their values rather than by pointer.
type LElement[T any] struct {
	e *Element[T]
	l *LList[T]
}

// Next returns the next list element or nil.
func (c *LElement[T]) Next() *LElement[T] {
	c.l.mu.RLock()
	defer c.l.mu.RUnlock()
	return c.l.wrap(c.e.Next())
}

// Prev returns the previous list element or nil.
func (c *LElement[T]) Prev() *LElement[T] {
	c.l.mu.RLock()
	defer c.l.mu.RUnlock()
	return c.l.wrap(c.e.Prev())
}

// Value returns the value stored with the element.
func (c *LElement[T]) Value() T {
	c.l.mu.RLock()
	defer c.l.mu.RUnlock()
	return c.e.Value
}

// Set replaces the value stored with the element.
func (c *LElement[T]) Set(v T) {
	c.l.mu.Lock()
	defer c.l.mu.Unlock()
	c.e.Value = v
}

// LList is a doubly linked list that is safe for concurrent use.
// Write operations hold mu exclusively, read operations hold it shared.
// The zero value for LList is an empty list ready to use.
//
// Elements are exposed as *LElement, whose methods use the list's lock.
type LList[T any] struct {
	list List[T]
	mu   sync.RWMutex
}

// Init initializes or clears list l.
func (l *LList[T]) Init() *LList[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.Init()
	return l
}

// NewL returns an initialized concurrent list.
func NewL[T any]() *LList[T] { return new(LList[T]).Init() }

// Len returns the number of elements of list l.
func (l *LList[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.list.Len()
}

// Front returns the first element of list l or nil if the list is empty.
func (l *LList[T]) Front() *LElement[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.wrap(l.list.Front())
}

// Back returns the last element of list l or nil if the list is empty.
func (l *LList[T]) Back() *LElement[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.wrap(l.list.Back())
}

func (l *LList[T]) wrap(e *Element[T]) *LElement[T] {
	if e == nil {
		return nil
	}
	return &LElement[T]{e: e, l: l}
}

// values returns a snapshot of the values of l, front to back.
func (l *LList[T]) values() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	vs := make([]T, 0, l.list.Len())
	for e := l.list.Front(); e != nil; e = e.Next() {
		vs = append(vs, e.Value)
	}
	return vs
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (l *LList[T]) Remove(e *LElement[T]) T {
	if e.l != l {
		return e.Value()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.list.Remove(e.e)
}

// PushFront inserts a new element with value v at the front of list l.
func (l *LList[T]) PushFront(v T) *LElement[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.wrap(l.list.PushFront(v))
}

// PushBack inserts a new element with value v at the back of list l.
func (l *LList[T]) PushBack(v T) *LElement[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.wrap(l.list.PushBack(v))
}

// InsertBefore inserts a new element with value v immediately before mark.
// If mark is not an element of l, the list is not modified.
func (l *LList[T]) InsertBefore(v T, mark *LElement[T]) *LElement[T] {
	if mark.l != l {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.wrap(l.list.InsertBefore(v, mark.e))
}

// InsertAfter inserts a new element with value v immediately after mark.
// If mark is not an element of l, the list is not modified.
func (l *LList[T]) InsertAfter(v T, mark *LElement[T]) *LElement[T] {
	if mark.l != l {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.wrap(l.list.InsertAfter(v, mark.e))
}

// MoveToFront moves element e to the front of list l.
func (l *LList[T]) MoveToFront(e *LElement[T]) {
	if e.l != l {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.MoveToFront(e.e)
}

// MoveToBack moves element e to the back of list l.
func (l *LList[T]) MoveToBack(e *LElement[T]) {
	if e.l != l {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.MoveToBack(e.e)
}

// MoveBefore moves element e to its new position before mark.
func (l *LList[T]) MoveBefore(e, mark *LElement[T]) {
	if e.l != l || mark.l != l {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.MoveBefore(e.e, mark.e)
}

// MoveAfter moves element e to its new position after mark.
func (l *LList[T]) MoveAfter(e, mark *LElement[T]) {
	if e.l != l || mark.l != l {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.MoveAfter(e.e, mark.e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The other list is snapshotted first, so the two locks are never held
// at the same time and l and other may be the same list.
func (l *LList[T]) PushBackList(other *LList[T]) {
	vs := other.values()
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, v := range vs {
		l.list.PushBack(v)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The other list is snapshotted first, so the two locks are never held
// at the same time and l and other may be the same list.
func (l *LList[T]) PushFrontList(other *LList[T]) {
	vs := other.values()
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := len(vs) - 1; i >= 0; i-- {
		l.list.PushFront(vs[i])
	}
}
