package queue

import "slices"

type Queue [T comparable]struct {
	data []T
	maxSize int
}


func NewQueue[T comparable](maxSize int) *Queue[T] {
	return &Queue[T]{maxSize: maxSize}
}


func (q *Queue[T]) Push(v T) {
	q.data = append(q.data, v)
	if q.maxSize > 0 && len(q.data) > q.maxSize {
		q.data = q.data[1:]
	}
}


func (q *Queue[T]) Pop() T {
	v := q.data[0]
	q.data = q.data[1:]
	return v
}


func (q *Queue[T]) Top() T {
	return q.data[0]
}


func (q *Queue[T]) Remove(v T) bool {
	ind := slices.Index(q.data, v)
	if ind == -1 {
		return false
	}
	q.data = slices.Delete(
		q.data,
		ind,
		ind + 1,
	)
	return true
}


func (q *Queue[T]) Size() int {
	return len(q.data)
}


func (q *Queue[T]) Contain(v T) bool {
	return slices.Contains(q.data, v)
}


func (q *Queue[T]) Clear() {
	q.data = slices.Delete(q.data, 0, len(q.data))
}