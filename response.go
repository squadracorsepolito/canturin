package main

type response[T entity] struct {
	changed bool
	undo    func() (T, error)
	redo    func() (T, error)
}

func newUnchangedResponse[T entity]() *response[T] {
	return &response[T]{changed: false}
}

func newResponse[T entity](undo, redo func() (T, error)) *response[T] {
	return &response[T]{changed: true, undo: undo, redo: redo}
}
