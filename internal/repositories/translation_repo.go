package repositories

import (
	"context"

	db "github.com/Ari-Pari/backend/internal/db/sqlc"
)

type TranslationRepository interface {
	CreateTranslations(ctx context.Context, translations []db.CreateTranslationsParams) error
}

type translationRepository struct {
	querier db.Querier
}

func (t translationRepository) CreateTranslations(ctx context.Context, translations []db.CreateTranslationsParams) error {

	_, err := t.querier.CreateTranslations(ctx, translations)

	return err
}
