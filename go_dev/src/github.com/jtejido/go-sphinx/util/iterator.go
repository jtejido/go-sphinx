package util

type Iterable[V any] interface {
	Iterator() Iterator[V]
}

type Iterator[V any] interface {
	HasNext() bool
	Next() V
}

type iteratorImpl[V any] struct {
	items []V
	index int
}

func NewIterator[V any](items []V) *iteratorImpl[V] {
	return &iteratorImpl[V]{items, 0}
}

// Next returns the next item in the collection.
func (i *iteratorImpl[V]) Next() V {
	item := i.items[i.index]
	i.index++
	return item
}

// HasNext return true if there are values to be read.
func (i *iteratorImpl[V]) HasNext() bool {
	return i.index < len(i.items)
}
