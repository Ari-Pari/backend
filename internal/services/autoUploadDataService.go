package services

import (
	"context"

	"github.com/Ari-Pari/backend/internal/repositories"
)

type AutoUploadDataService interface {
	CreateTranslations(ctx context.Context) (int64, error)
}

type autoUploadDataService struct {
	translationRepo repositories.TranslationRepository
}

type Translation struct {
	EngName string `json:"engName"`
	RuName  string `json:"ruName"`
	ArmName string `json:"armName"`
}
