package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/magiconair/properties/assert"
	"mini-e-wallet/domain"
	"regexp"
	"testing"
	"time"
)

func TestTokenRepository_GetByToken(t *testing.T) {
	var (
		query   string
		returns = domain.Tokens{
			Id:         "1",
			AccountID:  "testing-123",
			Token:      "27edda98-ac3c-48f0-b189-d78f883c923a",
			Expiration: time.Time{},
			CreatedAt:  time.Time{},
		}
	)
	sql, _, _ = sq.Select("account_id").From("tokens").Where(sq.And{
		sq.Eq{"token": "token"},
		sq.GtOrEq{"expiration": time.Now()},
	}).PlaceholderFormat(sq.Dollar).ToSql()

	returns = append(returns, domain.Users{
		Username: "fadli",
		Email:    "feedlyy@gmail.com",
	})

	testCases := []struct {
		name        string
		expectedErr bool
		context     context.Context
		doMockDB    func(mock sqlmock.Sqlmock)
		expected    []domain.Users
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"username", "email"}).
					AddRow("fadli", "feedlyy@gmail.com"))
			},
			expected: returns,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("internal server error"))
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err    error
				res    []domain.Users
				db     *sql.DB
				mock   sqlmock.Sqlmock
				mockDB *sqlx.DB
			)
			db, mock, err = sqlmock.New()
			if err != nil {
				panic(err)
			}
			defer db.Close()
			tc.doMockDB(mock)

			mockDB = sqlx.NewDb(db, "postgres")

			repoDB := NewTokenRepository(mockDB)
			res, err = repoDB.GetByToken(tc.context)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, res)
		})
	}
}
