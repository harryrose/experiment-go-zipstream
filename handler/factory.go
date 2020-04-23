package handler

import (
	"context"
	"io"
)

// FileHandlerFactory is a list of handler constructors
// When Construct is called on the factory, all handlers are tested in turn to see whether they
// can build a handler for the specific filename.  If none exists, then nil is returned
type FileHandlerFactory []Constructor

// Construct checks each constructor to see whether it will build a handler that can process the
// specified file.  If it can, then it's construct method is called and its return value returned.
// Otherwise, if nothing can process the file, then nil is returned
func (fh FileHandlerFactory) Construct(ctx context.Context, filename string) Handler {
	for _, f := range fh {
		if f.CanHandle(ctx, filename) {
			return f.Construct(ctx)
		}
	}
	return nil
}

type Constructor interface{
	CanHandle(ctx context.Context, filename string) bool
	Construct(ctx context.Context) Handler
}

type Handler interface {
	Handle(ctx context.Context, reader io.Reader) <-chan Item
}
