package services

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type LabelService struct {
	dt *pg.Datastore
}

// NewLabelService creates a new instance of LabelService with the provided datastore.
func NewLabelService(dt *pg.Datastore) *LabelService {
	return &LabelService{dt: dt}
}

// GetLabels retrieves all labels from the database along with their associated transaction counts.
// It also includes the rules that reference each label.
func (l *LabelService) GetLabels(ctx context.Context) ([]entity.Label, error) {
	labels, err := l.dt.QueryLabels(ctx)
	if err != nil {
		return []entity.Label{}, err
	}

	return labels, nil
}

// GetLabelByKeyValue retrieves a specific label by key and value.
func (l *LabelService) Get(ctx context.Context, key, value string) (*entity.Label, error) {
	labels, err := l.dt.QueryLabels(ctx, pg.LabelKeyValueQueryFilter(key, value))
	if err != nil {
		return nil, err
	}

	if len(labels) == 0 {
		return nil, nil
	}

	return &labels[0], nil
}

// IsExists checks if a label with the given key and value exists in the database.
func (l *LabelService) IsExists(ctx context.Context, key, value string) (bool, error) {
	label, err := l.Get(ctx, key, value)
	if err != nil {
		return false, err
	}

	return label != nil, nil
}

// Create creates a new label in the database if it doesn't already exist.
// If the label already exists, this method returns without error.
func (l *LabelService) Create(ctx context.Context, key, value string) (entity.Label, error) {
	label, err := l.Get(ctx, key, value)
	if err != nil {
		return entity.Label{}, err
	}

	if label != nil {
		return *label, nil
	}

	newLabel := entity.Label{
		Key:   key,
		Value: value,
	}

	if err := l.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		id, err := w.WriteLabel(ctx, newLabel)
		if err != nil {
			return err
		}
		newLabel.ID = id
		return nil
	}); err != nil {
		return entity.Label{}, err
	}

	return newLabel, nil
}

// Delete removes a label from the database by its key and value.
func (l *LabelService) Delete(ctx context.Context, key, value string) error {
	label := entity.Label{
		Key:   key,
		Value: value,
	}

	return l.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		return w.DeleteLabel(ctx, label)
	})
}

// DeleteLabel removes a label entity from the database.
func (l *LabelService) DeleteLabel(ctx context.Context, label entity.Label) error {
	return l.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		return w.DeleteLabel(ctx, label)
	})
}
