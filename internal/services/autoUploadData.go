package services

import (
	"context"

	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/domain"
)

func (a *autoUploadDataService) CreateTranslations(ctx context.Context, translations []domain.Translation) error {
	translationsParams := make([]db.CreateTranslationsParams, len(translations))
	for i, translation := range translations {
		translationsParams[i] = ToDomain(translation)
	}
	return a.translationRepo.CreateTranslations(ctx, translationsParams)
}
