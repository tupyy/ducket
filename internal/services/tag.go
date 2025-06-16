package services

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type TagService struct {
	dt *pg.Datastore
}

func NewTagService(dt *pg.Datastore) *TagService {
	return &TagService{dt: dt}
}

func (t *TagService) GetTags(ctx context.Context) ([]entity.Tag, error) {
	return t.dt.QueryTags(ctx)
}

func (t *TagService) IsExists(ctx context.Context, tag string) (bool, error) {
	tags, err := t.dt.QueryTags(ctx)
	if err != nil {
		return false, err
	}

	for _, t := range tags {
		if t.Value == tag {
			return true, nil
		}
	}

	return false, nil
}

func (t *TagService) Create(ctx context.Context, tag string) error {
	exists, err := t.IsExists(ctx, tag)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return t.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.WriteTag(ctx, tag)
	})
}

func (t *TagService) Delete(ctx context.Context, tag string) error {
	return t.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.DeleteTag(ctx, tag)
	})
}
