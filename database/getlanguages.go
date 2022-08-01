package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type (
	// Language it is used for getting records from DB
	Language struct {
		Id           int    `json:"id"`
		Country      string `json:"country"`
		Language     string `json:"language"`
		TwoLetters   string `json:"two_letters"`
		ThreeLetters string `json:"three_letters"`
		Number       int64  `json:"number"`
	}
)

// GetLanguages Used to get all Languages from the languages table
func GetLanguages(session *sql.DB) ([]byte, error) {

	var languages []Language

	// Query to get all language records
	rows, err := session.Query("SELECT * FROM languages")
	if err != nil {
		return nil, fmt.Errorf("unable to execute SQL query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lng Language
		if err = rows.Scan(&lng.Id, &lng.Country, &lng.Language, &lng.TwoLetters, &lng.ThreeLetters, &lng.Number); err != nil {
			return nil, fmt.Errorf("unable to collect query result: %w", err)
		}
		languages = append(languages, lng)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	data, err := json.Marshal(languages)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal data: %w", err)
	}

	return []byte(fmt.Sprintf("{\"Body\":[%s]}", string(data))), nil
}
