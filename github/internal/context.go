package internal

import "context"

type ContextAccessor func() context.Context

func (c ContextAccessor) Context() context.Context {
	return c()
}

func BackgroundContext() ContextAccessor {
	return context.Background
}

func TodoContext() ContextAccessor {
	return context.TODO
}
