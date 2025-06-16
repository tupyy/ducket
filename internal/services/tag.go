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
	// count transaction first
	stats, err := t.dt.CountTransactions(ctx)
	if err != nil {
		return []entity.Tag{}, err
	}

	tags, err := t.dt.QueryTags(ctx)
	if err != nil {
		return []entity.Tag{}, err
	}

	modifiedTags := make(map[string]entity.Tag)
	for _, tag := range tags {
		mTag, ok := modifiedTags[tag.Value]
		if !ok {
			mTag = entity.Tag{
				Value:     tag.Value,
				Rules:     tag.Rules,
				CreatedAt: tag.CreatedAt,
			}
		}
		for _, s := range stats {
			if s.Tag == tag.Value {
				mTag.CountTransactions += s.Count
			}
		}
		modifiedTags[mTag.Value] = mTag
	}

	exportedTags := make([]entity.Tag, 0, len(modifiedTags))
	for _, t := range modifiedTags {
		exportedTags = append(exportedTags, t)
	}

	return exportedTags, nil
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
