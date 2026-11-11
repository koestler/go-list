package list

import (
	"sync"
	"testing"
)

func checkLListPointers[T any](t *testing.T, l *LList[T], es []*LElement[T]) {
	raw := make([]*Element[T], len(es))
	for i, e := range es {
		raw[i] = e.e
	}
	checkListPointers(t, &l.list, raw)
}

func TestLList(t *testing.T) {
	l := NewL[any]()
	checkLListPointers(t, l, []*LElement[any]{})

	// Single element list
	e := l.PushFront("a")
	checkLListPointers(t, l, []*LElement[any]{e})
	l.MoveToFront(e)
	checkLListPointers(t, l, []*LElement[any]{e})
	l.MoveToBack(e)
	checkLListPointers(t, l, []*LElement[any]{e})
	l.Remove(e)
	checkLListPointers(t, l, []*LElement[any]{})

	// Bigger list
	e2 := l.PushFront(2)
	e1 := l.PushFront(1)
	e3 := l.PushBack(3)
	e4 := l.PushBack("banana")
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e3, e4})

	l.Remove(e2)
	checkLListPointers(t, l, []*LElement[any]{e1, e3, e4})

	l.MoveToFront(e3) // move from middle
	checkLListPointers(t, l, []*LElement[any]{e3, e1, e4})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e3})

	l.MoveToFront(e3) // move from back
	checkLListPointers(t, l, []*LElement[any]{e3, e1, e4})
	l.MoveToFront(e3) // should be no-op
	checkLListPointers(t, l, []*LElement[any]{e3, e1, e4})

	l.MoveToBack(e3) // move from front
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e3})
	l.MoveToBack(e3) // should be no-op
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e3})

	e2 = l.InsertBefore(2, e1) // insert before front
	checkLListPointers(t, l, []*LElement[any]{e2, e1, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(2, e4) // insert before middle
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(2, e3) // insert before back
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e2, e3})
	l.Remove(e2)

	e2 = l.InsertAfter(2, e1) // insert after front
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(2, e4) // insert after middle
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e2, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(2, e3) // insert after back
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e3, e2})
	l.Remove(e2)

	// Check standard iteration.
	sum := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value().(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all elements by iterating
	var next *LElement[any]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
	checkLListPointers(t, l, []*LElement[any]{})
}

func checkLList(t *testing.T, l *LList[any], es []any) {
	if !checkListLen(t, &l.list, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le := e.Value().(int)
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func TestLExtending(t *testing.T) {
	l1 := NewL[any]()
	l2 := NewL[any]()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := NewL[any]()
	l3.PushBackList(l1)
	checkLList(t, l3, []any{1, 2, 3})
	l3.PushBackList(l2)
	checkLList(t, l3, []any{1, 2, 3, 4, 5})

	l3 = NewL[any]()
	l3.PushFrontList(l2)
	checkLList(t, l3, []any{4, 5})
	l3.PushFrontList(l1)
	checkLList(t, l3, []any{1, 2, 3, 4, 5})

	checkLList(t, l1, []any{1, 2, 3})
	checkLList(t, l2, []any{4, 5})

	l3 = NewL[any]()
	l3.PushBackList(l1)
	checkLList(t, l3, []any{1, 2, 3})
	l3.PushBackList(l3)
	checkLList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = NewL[any]()
	l3.PushFrontList(l1)
	checkLList(t, l3, []any{1, 2, 3})
	l3.PushFrontList(l3)
	checkLList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = NewL[any]()
	l1.PushBackList(l3)
	checkLList(t, l1, []any{1, 2, 3})
	l1.PushFrontList(l3)
	checkLList(t, l1, []any{1, 2, 3})
}

func TestLRemove(t *testing.T) {
	l := NewL[any]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	checkLListPointers(t, l, []*LElement[any]{e1, e2})
	e := l.Front()
	l.Remove(e)
	checkLListPointers(t, l, []*LElement[any]{e2})
	l.Remove(e)
	checkLListPointers(t, l, []*LElement[any]{e2})
}

func TestLIssue4103(t *testing.T) {
	l1 := NewL[any]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := NewL[any]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestLIssue6349(t *testing.T) {
	l := NewL[any]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e)
	if e.Value() != 1 {
		t.Errorf("e.value = %d, want 1", e.Value())
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestLMove(t *testing.T) {
	l := NewL[any]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	l.MoveAfter(e3, e3)
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e3, e4})

	l.MoveAfter(e3, e2)
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e3, e4})
	l.MoveBefore(e2, e3)
	checkLListPointers(t, l, []*LElement[any]{e1, e2, e3, e4})

	l.MoveBefore(e2, e4)
	checkLListPointers(t, l, []*LElement[any]{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.MoveBefore(e4, e1)
	checkLListPointers(t, l, []*LElement[any]{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkLListPointers(t, l, []*LElement[any]{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e2, e3)
	checkLListPointers(t, l, []*LElement[any]{e1, e3, e2, e4})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
func TestLZeroList(t *testing.T) {
	var l1 = new(LList[any])
	l1.PushFront(1)
	checkLList(t, l1, []any{1})

	var l2 = new(LList[any])
	l2.PushBack(1)
	checkLList(t, l2, []any{1})

	var l3 = new(LList[any])
	l3.PushFrontList(l1)
	checkLList(t, l3, []any{1})

	var l4 = new(LList[any])
	l4.PushBackList(l2)
	checkLList(t, l4, []any{1})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestLInsertBeforeUnknownMark(t *testing.T) {
	var l LList[any]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertBefore(1, l.wrap(new(Element[any])))
	checkLList(t, &l, []any{1, 2, 3})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestLInsertAfterUnknownMark(t *testing.T) {
	var l LList[any]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertAfter(1, l.wrap(new(Element[any])))
	checkLList(t, &l, []any{1, 2, 3})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestLMoveUnknownMark(t *testing.T) {
	var l1 LList[any]
	e1 := l1.PushBack(1)

	var l2 LList[any]
	e2 := l2.PushBack(2)

	l1.MoveAfter(e1, e2)
	checkLList(t, &l1, []any{1})
	checkLList(t, &l2, []any{2})

	l1.MoveBefore(e1, e2)
	checkLList(t, &l1, []any{1})
	checkLList(t, &l2, []any{2})
}

func TestLListConcurrent(t *testing.T) {
	var l LList[int]
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Go(func() {
			for j := 0; j < 200; j++ {
				e := l.PushBack(j)
				e.Set(j + 1)
				l.MoveToFront(e)
				for x := l.Front(); x != nil; x = x.Next() {
					_ = x.Value()
				}
				l.Remove(e)
			}
		})
	}
	wg.Wait()
	if l.Len() != 0 {
		t.Fatalf("len = %d, want 0", l.Len())
	}
}

// Test that operations on an element of another list neither modify the
// lists nor race with concurrent changes to the other list.
func TestLListForeignElementConcurrent(t *testing.T) {
	l1, l2 := NewL[int](), NewL[int]()
	foreign := l2.PushBack(1)
	l1.PushBack(1)

	var wg sync.WaitGroup
	wg.Go(func() {
		for i := range 500 {
			foreign.Set(i)
			l2.MoveToBack(foreign)
		}
	})
	wg.Go(func() {
		for range 500 {
			l1.Remove(foreign)
			l1.MoveToFront(foreign)
			l1.MoveToBack(foreign)
			l1.MoveBefore(foreign, l1.Front())
			l1.MoveAfter(foreign, l1.Front())
			if l1.InsertBefore(2, foreign) != nil || l1.InsertAfter(2, foreign) != nil {
				t.Error("insert with foreign mark must return nil")
			}
		}
	})
	wg.Wait()

	if n := l1.Len(); n != 1 {
		t.Errorf("l1.Len() = %d, want 1", n)
	}
	if n := l2.Len(); n != 1 {
		t.Errorf("l2.Len() = %d, want 1", n)
	}
}
