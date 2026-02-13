package services

import (
	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToDomain(translation domain.Translation) db.CreateTranslationsParams {
	return db.CreateTranslationsParams{
		EngName: pgtype.Text{String: translation.EngName, Valid: true},
		RuName:  pgtype.Text{String: translation.RuName, Valid: true},
		ArmName: pgtype.Text{String: translation.ArmName, Valid: true},
	}
}
