package api

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/Ari-Pari/backend/internal/api/generated"
	"github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// Mock для БД (Querier)
type mockDB struct {
	db.Querier
	GetDanceByIDFunc func(params db.GetDanceByIDParams) (db.GetDanceByIDRow, error)
}

func (m *mockDB) GetDanceByID(ctx context.Context, p db.GetDanceByIDParams) (db.GetDanceByIDRow, error) {
	return m.GetDanceByIDFunc(p)
}

// Mock для MinIO
type mockStorage struct {
	GetFileURLFunc func(key string) (string, error)
}

func (m *mockStorage) GetFileURL(ctx context.Context, key string, exp time.Duration) (string, error) {
	return m.GetFileURLFunc(key)
}

func (m *mockStorage) UploadImage(context.Context, string, io.Reader, int64, string) (string, error) { return "", nil }
func (m *mockStorage) DeleteFile(context.Context, string) error { return nil }
func (m *mockStorage) GetOriginalName(context.Context, string) (string, error) { return "", nil }



func TestGetDancesId(t *testing.T) {
	tests := []struct {
		name           string
		danceID        int
		mockDBResp     db.GetDanceByIDRow
		mockDBErr      error
		expectedStatus int
	}{
		{
			name:    "Success 200",
			danceID: 1,
			mockDBResp: db.GetDanceByIDRow{
				ID:         1,
				Name:       "Berd",
				Complexity: pgtype.Int4{Int32: 3, Valid: true},
				PhotoKey:   pgtype.Text{String: "photo.jpg", Valid: true},
				Gender:     "male",
				Paces:      []int32{1, 2},
				RegionsJson: []byte(`[{"id": 10, "name": "Shirak"}]`),
			},
			mockDBErr:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found 404",
			danceID:        999,
			mockDBErr:      pgx.ErrNoRows,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Internal Error 500",
			danceID:        1,
			mockDBErr:      errors.New("db connection lost"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mDB := &mockDB{
				GetDanceByIDFunc: func(p db.GetDanceByIDParams) (db.GetDanceByIDRow, error) {
					return tt.mockDBResp, tt.mockDBErr
				},
			}
			mStorage := &mockStorage{
				GetFileURLFunc: func(key string) (string, error) {
					return "http://minio/" + key, nil
				},
			}

			logger := log.New(io.Discard, "", 0) 
			srv := NewServer(logger, mDB, mStorage)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/dances/1", nil)
			w := httptest.NewRecorder()

			srv.GetDancesId(w, req, tt.danceID, api.GetDancesIdParams{})

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, w.Body.String(), tt.mockDBResp.Name)
				assert.Contains(t, w.Body.String(), "Shirak")
			}
		})
	}
}