package handlers

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
)

// MustFromContext retrieves the store from the context or panics if not found.
// This function is used with gin middleware to extract the store that was
// injected by the datastore middleware.
func MustFromContext(ctx context.Context) *store.Store {
	// this is for gin middleware which does not accept key as any. only string.
	if c := ctx.Value("datastore"); c != nil {
		return c.(*store.Store)
	}
	panic("datastore middleware did not inject store")
}
