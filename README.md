# go-list
Generic linked lists in go.

This library copied the standard library package `container/list` and made the `List` and `Element` types generic.
It provides two implementations:

| Type | Element | Concurrency |
|------|---------|-------------|
| `List[T]` | `Element[T]` | not safe for concurrent use |
| `LList[T]` | `LElement[T]` | safe for concurrent use |

## List

`List[T]` is a drop-in generic version of `container/list.List`. The zero value is an empty list ready to use.
Unlike `container/list`, `Init` detaches all existing elements (O(n)), so stale elements can no longer corrupt the list.

```go
l := list.New[string]()
l.PushBack("b")
l.PushFront("a")

for e := l.Front(); e != nil; e = e.Next() {
	fmt.Println(e.Value) // a, b
}
```

## LList

`LList[T]` (locked list) wraps `List[T]` and guards it with a `sync.RWMutex`: write operations take the write lock,
read operations the read lock. It offers the same methods as `List[T]`, but works with `*LElement[T]`, whose
`Next`, `Prev`, `Value` and `Set` methods use the list's lock as well. Use `Value()` and `Set(v)` instead of the
`Value` field. The zero value is an empty list ready to use.

```go
l := list.NewL[int]()

var wg sync.WaitGroup
for i := 0; i < 4; i++ {
	wg.Go(func() {
		e := l.PushBack(i)
		e.Set(e.Value() * 10)
	})
}
wg.Wait()

for e := l.Front(); e != nil; e = e.Next() {
	fmt.Println(e.Value())
}
```

Notes:
* Each `LList` method returns a new `*LElement` wrapper, so compare elements by value, not by pointer.
* Passing an element that belongs to another list is a no-op, as with `List`.
* `PushBackList` and `PushFrontList` copy a snapshot of the other list; the list is not locked as a whole across
  multiple calls. Iterating with `Next` is safe but not atomic: other goroutines may modify the list in between.
