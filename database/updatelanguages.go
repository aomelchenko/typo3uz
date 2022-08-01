package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const SaveError = "Save failed."

// Changes used to assign POST variable of POST:TGData[Changes]
type Changes struct {
	Id           string `json:"id"`
	Added        int    `json:"added"`
	Changed      int    `json:"changed"`
	Deleted      int    `json:"deleted"`
	Language     string `json:"language"`
	Country      string `json:"country"`
	TwoLetters   string `json:"two_letters"`
	ThreeLetters string `json:"three_letters"`
	Number       string `json:"number"`
}

// PostRequest used to assign POST request POST:TGData
type PostRequest struct {
	IO struct {
		Message string
	}
	Changes []Changes
}

// Response describes response data structure
type Response struct {
	IO struct {
		HtmlMessage string
		Result      int
	}
	Changes []ChangedRow
}

// ChangedRow used to return Messages for POST update
type ChangedRow struct {
	Id      string `json:"id"`
	Changed int
	Added   int
	Deleted int
	Color   string
}

// UpdateLanguages method handles update the languages table requests
func UpdateLanguages(req *http.Request, session *sql.DB) ([]byte, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	jsonString := req.Form.Get("Data")

	var post PostRequest

	// Convert string to JSON format and assign it to post variable
	err = json.Unmarshal([]byte(jsonString), &post)
	if err != nil {
		return nil, err
	}

	resp := Response{}

	allChanges, err := processChanges(session, post)
	if err != nil {
		if strings.Contains(err.Error(), SaveError) {
			resp.IO.HtmlMessage = err.Error()
			resp.Changes = allChanges
			ret, err1 := json.Marshal(resp)
			if err1 != nil {
				return nil, err1
			}
			return ret, nil
		}
		return nil, err
	}

	resp.Changes = allChanges
	ret, err := json.Marshal(resp)

	if err != nil {
		return nil, err
	}
	return ret, nil
}

func rowCreate(session *sql.DB, row Changes) error {
	err := checkRowExist(session, row)
	if err != nil {
		return err
	}

	// Prepare row Insert query
	stmt, err := session.Prepare("INSERT INTO  languages (country, language, two_letters, three_letters, number) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Assign insert variables - this way it's more secure
	_, err = stmt.Exec(row.Country, row.Language, row.TwoLetters, row.ThreeLetters, row.Number)
	if err != nil {
		return err
	}
	return nil
}

func rowUpdate(session *sql.DB, row Changes) error {
	err := checkRowExist(session, row)
	if err != nil {
		return err
	}

	// Update is made on various fields, sometimes it's one, sometimes there are many, so lets check which fields should be updated and later use strings.Join similar to PHP array implode
	var s []string
	var v []interface{}

	// fields to be excluded from writing to db
	excludeFields := map[string]bool{
		"id":      true,
		"added":   true,
		"deleted": true,
		"changed": true,
	}

	var inInterface map[string]interface{}
	inrec, err := json.Marshal(row)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}

	// iterate through inrecs
	for field, val := range inInterface {
		if !excludeFields[field] && val != "" {
			s = append(s, strings.ToLower(field)+" = ?")
			v = append(v, val)
		}
	}

	v = append(v, row.Id)
	// Prepare update statement
	stmt, err := session.Prepare("UPDATE  languages SET " + strings.Join(s, ", ") + " WHERE id=?")
	if err != nil {
		return err
	}

	// v... is used for slicing the array of interfaces and pass them as arguments
	_, err = stmt.Exec(v...)
	if err != nil {
		return err
	}

	return nil
}

func rowDelete(session *sql.DB, row Changes) error {
	// Prepare record delete query
	stmt, err := session.Prepare("DELETE from  languages where id=?")
	if err != nil {
		log.Fatalln(err)
	}
	// Close the resource
	defer stmt.Close()

	_, err = stmt.Exec(row.Id)
	if err != nil {
		return err
	}

	return nil
}

func processChanges(session *sql.DB, post PostRequest) ([]ChangedRow, error) {
	var allChanges []ChangedRow

	// Iterate through all Post.Changes
	for _, row := range post.Changes {
		changedRow := ChangedRow{}
		changedRow.Id = row.Id
		changedRow.Color = "rgb(255, 255, 166)"

		if row.Added == 1 {
			err := rowCreate(session, row)
			changedRow.Added = 1
			if err != nil {
				changedRow.Added = -1
				changedRow.Color = "rgb(255, 0, 0)"
				return append(allChanges, changedRow), err
			}
			allChanges = append(allChanges, changedRow)
		}

		if row.Deleted == 1 {
			err := rowDelete(session, row)
			changedRow.Deleted = 1
			if err != nil {
				changedRow.Deleted = -1
				changedRow.Color = "rgb(255, 0, 0)"
				return append(allChanges, changedRow), err
			}
			allChanges = append(allChanges, changedRow)
		}

		if row.Changed == 1 {
			err := rowUpdate(session, row)
			changedRow.Changed = 1
			if err != nil {
				changedRow.Changed = -1
				changedRow.Color = "rgb(255, 0, 0)"
				return append(allChanges, changedRow), err
			}
			allChanges = append(allChanges, changedRow)
		}
	}
	return allChanges, nil
}

func checkRowExist(session *sql.DB, row Changes) error {
	if len(row.Country) == 0 || len(row.Number) == 0 {
		// Cover a partial row update case
		var country, number string
		rows, err := session.Query(fmt.Sprintf("SELECT country, number FROM languages WHERE id = %s", row.Id))
		if err != nil {
			return err
		}
		defer rows.Close()

		if rows.Next() {
			if err = rows.Scan(&country, &number); err != nil {
				return fmt.Errorf("unable to collect query result: %w", err)
			}
		}
		// Add missing values
		if len(row.Country) == 0 {
			row.Country = country
		}
		if len(row.Number) == 0 {
			row.Number = number
		}
	}

	// Check if the row with this name and number exists
	rows, err := session.Query(fmt.Sprintf("SELECT * FROM languages WHERE country = '%s' AND number = %s", row.Country, row.Number))
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return fmt.Errorf(fmt.Sprintf("%s Record with number `%s` already exist", SaveError, row.Number))
	}
	return nil
}
