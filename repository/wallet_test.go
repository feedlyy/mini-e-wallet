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

func TestWalletRepository_GetByOwnedID(t *testing.T) {
	var (
		query   string
		returns = domain.Wallets{
			Id:         "1",
			OwnedBy:    "testing-123",
			Status:     "enabled",
			EnableAt:   time.Time{},
			Balance:    20000,
			DisabledAt: sql.NullTime{},
		}
		enableAt, disabledAt time.Time
		id                   = "27edda98-ac3c-48f0-b189-d78f883c923a"
	)
	query, _, _ = sq.Select("id", "owned_by", "status", "enabled_at", "balance", "disabled_at").From("wallets").Where(sq.And{
		sq.Eq{"owned_by": "id"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	enableAt, _ = time.Parse("2023-04-01 11:06:55.696913", "2023-04-01 11:06:55.696913")
	disabledAt, _ = time.Parse("2023-04-01 11:06:55.696913", "2023-04-01 11:06:55.696913")
	returns.EnableAt = enableAt
	returns.DisabledAt = sql.NullTime{
		Time:  disabledAt,
		Valid: true,
	}

	testCases := []struct {
		name        string
		expectedErr bool
		context     context.Context
		doMockDB    func(mock sqlmock.Sqlmock)
		expected    domain.Wallets
		input       string
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "owned_by", "status", "enabled_at", "balance", "disabled_at"}).
						AddRow("1", "testing-123", "enabled", enableAt, 20000, disabledAt))
			},
			input:    id,
			expected: returns,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs(id).WillReturnError(errors.New("internal server error"))
			},
			expected: domain.Wallets{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err    error
				res    domain.Wallets
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

			repoDB := NewWalletRepository(mockDB)
			res, err = repoDB.GetByOwnedID(tc.context, tc.input)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, res)
		})
	}
}
