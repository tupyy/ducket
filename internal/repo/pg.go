package repo

import (
	"context"
	"errors"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/repo/models"
	"gorm.io/gorm"
)

type TransationRepo struct {
	db             *gorm.DB
	client         Client
	circuitBreaker CircuitBreaker
}

func NewRepo(client Client) (*TransationRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &TransationRepo{}, err
	}

	return &TransationRepo{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (t *TransationRepo) Write(ctx context.Context, transaction *entity.Transaction) error {
	model := models.Transaction{
		CreatedAt:       time.Now(),
		Date:            transaction.Date,
		TransactionType: string(transaction.Kind),
		Description:     transaction.RawContent,
		Amount:          float64(transaction.Sum),
	}

	tx := t.db.WithContext(ctx).Begin()

	if err := tx.Create(&model).Error; err != nil {
		tx.Rollback()
		return err
	}

	// start creating labels
	for key, value := range transaction.Labels {
		label, err := t.getLabel(tx, key, value)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				l, err := t.addLabel(tx, key, value)
				if err != nil {
					tx.Rollback()
					return err
				}
				label = l
			}
		}
		if err := t.associateTransactionWithLabel(tx, label, model); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (t *TransationRepo) getLabel(tx *gorm.DB, key, value string) (models.Label, error) {
	label := models.Label{}
	if err := tx.Where("key = ? AND value = ?", key, value).First(&label).Error; err != nil {
		return label, err
	}
	return label, nil
}

func (t *TransationRepo) addLabel(tx *gorm.DB, key, value string) (models.Label, error) {
	label := models.Label{
		Key:   key,
		Value: value,
	}
	if err := tx.Create(&label).Error; err != nil {
		return label, err
	}
	return label, nil
}

func (t *TransationRepo) associateTransactionWithLabel(tx *gorm.DB, label models.Label, transaction models.Transaction) error {
	a := models.TransactionsLabels{
		TransactionID: transaction.ID,
		LabelID:       label.ID,
	}
	return tx.Create(&a).Error
}
