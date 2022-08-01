package database

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkGetLanguages(b *testing.B) {
	session := InitDb()
	for i := 0; i < 100; i++ {
		_, _ = GetLanguages(session)
	}
}

func TestGetLanguages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "country", "language", "two_letters", "three_letters", "number"}).
		AddRow(1, "Ukraine", "Ukrainian", "uk-UA", "uk-UKR", "1991")
	mock.ExpectQuery("^SELECT (.+) FROM languages$").WillReturnRows(rows)

	languages, err := GetLanguages(db)
	require.NoError(t, err)
	require.Contains(t, string(languages), "Ukrainian")

	testError := errors.New("some error")
	mock.ExpectQuery("^SELECT (.+) FROM languages$").WillReturnError(testError)
	languages, err = GetLanguages(db)
	require.ErrorIs(t, err, testError)
	require.Empty(t, languages)
}
