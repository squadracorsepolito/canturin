package main

type response[T entity] struct {
	changed bool
	undo    func() (T, error)
	redo    func() (T, error)
}

func newResponse[T entity]() *response[T] {
	return &response[T]{changed: false}
}

func (r *response[T]) setUndo(undo func() (T, error)) {
	r.undo = undo
	r.changed = true
}

func (r *response[T]) setRedo(redo func() (T, error)) {
	r.redo = redo
	r.changed = true
}
