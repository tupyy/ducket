package handlers

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
)

func MustFromContext(ctx context.Context) *pg.Datastore {
	// this is for gin middleware which does not accept key as any. only string.
	if c := ctx.Value("datastore"); c != nil {
		return c.(*pg.Datastore)
	}
	panic("datastore middleware did not inject datastore")
}
