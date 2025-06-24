package services

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type TagService struct {
	dt *pg.Datastore
}

// NewTagService creates a new instance of TagService with the provided datastore.
func NewTagService(dt *pg.Datastore) *TagService {
	return &TagService{dt: dt}
}

// GetTags retrieves all tags from the database along with their associated transaction counts.
// It also includes the rules that reference each tag.
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

	// Add number of transactions for each tag
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

// IsExists checks if a tag with the given value exists in the database.
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

// Create creates a new tag in the database if it doesn't already exist.
// If the tag already exists, this method returns without error.
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

// Delete removes a tag from the database by its value.
func (t *TagService) Delete(ctx context.Context, tag string) error {
	return t.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.DeleteTag(ctx, tag)
	})
}
