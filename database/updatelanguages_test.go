package database

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckRowExist(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		changes  Changes
		query    string
		rows     *sqlmock.Rows
		subQuery string
		subRows  *sqlmock.Rows
		wantErr  bool
		err      error
	}{
		{
			name: "both values present",
			changes: Changes{
				Id:      "1",
				Country: "Ukraine",
				Number:  "1991",
			},
			query: "^SELECT (.+) FROM languages WHERE country = 'Ukraine' AND number = 1991$",
			rows: sqlmock.NewRows([]string{"id", "country", "language", "two_letters", "three_letters", "number"}).
				AddRow(1, "Ukraine", "Ukrainian", "uk-UA", "uk-UKR", "1991"),
			subQuery: "",
			subRows:  nil,
			wantErr:  true,
			err:      nil,
		}, {
			name: "no records",
			changes: Changes{
				Id:      "1",
				Country: "Ukraine",
				Number:  "1992",
			},
			query:    "^SELECT (.+) FROM languages WHERE country = 'Ukraine' AND number = 1992$",
			rows:     sqlmock.NewRows(nil),
			subQuery: "",
			subRows:  nil,
			wantErr:  false,
			err:      nil,
		}, {
			name: "one value present",
			changes: Changes{
				Id:     "1",
				Number: "1991",
			},
			query: "^SELECT (.+) FROM languages WHERE country = 'Ukraine' AND number = 1991$",
			rows: sqlmock.NewRows([]string{"id", "country", "language", "two_letters", "three_letters", "number"}).
				AddRow(1, "Ukraine", "Ukrainian", "uk-UA", "uk-UKR", "1991"),
			subQuery: "^SELECT country, number FROM languages WHERE id = 1$",
			subRows: sqlmock.NewRows([]string{"country", "number"}).
				AddRow("Ukraine", "1991"),
			wantErr: true,
			err:     nil,
		}, {
			name: "db error present",
			changes: Changes{
				Id:      "1",
				Country: "Ukraine",
				Number:  "1991",
			},
			query:   "^SELECT (.+) FROM languages WHERE country = 'Ukraine' AND number = 1991$",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.err == nil {
				if tt.subRows != nil {
					mock.ExpectQuery(tt.subQuery).WillReturnRows(tt.subRows)
				}
				mock.ExpectQuery(tt.query).WillReturnRows(tt.rows)
			} else {
				mock.ExpectQuery(tt.query).WillReturnError(tt.err)
			}
			err = checkRowExist(db, tt.changes)
			if tt.wantErr {
				require.Error(t, err)
				if tt.err != nil {
					require.ErrorIs(t, err, tt.err, "unexpected error")
				} else {
					require.Contains(t, err.Error(), SaveError)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
