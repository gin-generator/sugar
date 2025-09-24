package bootstrap

import "context"

type Option interface {
	apply(*Bootstrap)
}

type optionFunc func(*Bootstrap)

func (o optionFunc) apply(b *Bootstrap) {
	o(b)
}

type Bootstrap struct {
	context.Context

	// other server options
}
