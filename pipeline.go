package pipeline

// A Handler that let's you build a processing unit
// for your data pipeline
type Pipe[T any] func(p T) T

// Empty pipe that passes the data it recieves
// without modifying it
func EmptyPipe[T any]() Pipe[T] {
	return func(p T) T {
		return p
	}
}

// Pipes the data passed to it
// through all the data handlers (pipes)
// containted within
type Pipeline[T any] struct {
	pipes []Pipe[T]
}

// Returns a new pipeline with the given pipes in order
func New[T any](pipes ...Pipe[T]) Pipeline[T] {
	return Pipeline[T]{pipes}
}

// Empty pipeline without any pipes
func Empty[T any]() Pipeline[T] {
	return New[T]()
}

// Closes the pipeline and returns the final value
// NOTE: Use when the the pipes are handling a value,
// that needs to be passed in order
func (p Pipeline[T]) CloseWith(val T) T {
	for i := range p.pipes {
		val = p.pipes[i](val)
	}
	return val
}

// Closes the pipeline and returns the final function/interface
// NOTE: Use when the the pipes are building a function,
// that will be called later
func (p Pipeline[T]) EmbedFinal(fn T) T {
	for i := range p.pipes {
		fn = p.pipes[len(p.pipes)-1-i](fn)
	}
	return fn
}

// Pipe the data from a pipeline into raw handlers
func (p Pipeline[T]) IntoRaw(funcs ...Pipe[T]) Pipeline[T] {
	list := make([]Pipe[T], 0, len(p.pipes)+len(funcs))
	list = append(list, p.pipes...)
	list = append(list, funcs...)
	return New(list...)
}

// Pipe the data from a pipeline into another pipeline
func (p Pipeline[T]) Into(pl Pipeline[T]) Pipeline[T] {
	return p.IntoRaw(pl.pipes...)
}
