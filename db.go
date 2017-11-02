package csnotes

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// SetUpDB creates all the tables necessary for the app.
func SetUpDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE users (
				id		INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				name	VARCHAR(191) NOT NULL UNIQUE,
				admin	BOOLEAN DEFAULT FALSE NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`CREATE TABLE notes (
				id		INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				title	VARCHAR(191) NOT NULL,
				content TEXT,
				time	DATETIME,
				user_id	INT(10) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE tags (
				id		INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				title	VARCHAR(191) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE note_tag (
				note_id	INT(10) NOT NULL,
				tag_id	INT(10) NOT NULL,
				PRIMARY KEY (note_id, tag_id)
			)`)
	if err != nil {
		return err
	}

	return nil
}

// TearDownDB clears the database of all tables that the app uses.
func TearDownDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS notes")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS tags")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS note_tag")
	if err != nil {
		return err
	}

	return nil
}

func SeedDB(db *sql.DB) (ids map[string]int64, err error) {
	ids = map[string]int64{}

	// Users
	res, err := db.Exec("INSERT INTO users (name, admin) VALUES (?, ?)", "nonadmin", false)
	if err != nil {
		return
	}
	ids["user.nonadmin"], err = res.LastInsertId()
	if err != nil {
		return
	}

	res, err = db.Exec("INSERT INTO users (name, admin) VALUES (?, ?)", "admin", true)
	if err != nil {
		return
	}
	ids["user.admin"], err = res.LastInsertId()
	if err != nil {
		return
	}

	// Notes
	res, err = db.Exec("INSERT INTO notes (title, content, time, user_id) VALUES (?, ?, ?, ?)",
		"note1", "content", "2017-01-01 12:00", ids["user.nonadmin"])
	if err != nil {
		return
	}
	ids["note.note1"], err = res.LastInsertId()
	if err != nil {
		return
	}

	res, err = db.Exec("INSERT INTO notes (title, content, time, user_id) VALUES (?, ?, ?, ?)",
		"note2", "content", "2017-02-01 12:00", ids["user.nonadmin"])
	if err != nil {
		return
	}
	ids["note.note2"], err = res.LastInsertId()
	if err != nil {
		return
	}

	// Tags
	res, err = db.Exec("INSERT INTO tags (title) VALUES (?)", "tag1")
	if err != nil {
		return
	}
	ids["tag.tag1"], err = res.LastInsertId()
	if err != nil {
		return
	}

	res, err = db.Exec("INSERT INTO tags (title) VALUES (?)", "tag2")
	if err != nil {
		return
	}
	ids["tag.tag2"], err = res.LastInsertId()
	if err != nil {
		return
	}

	// Attach tags to notes. Note 1 will be tagged with tag 1, etc.
	_, err = db.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)",
		ids["note.note1"], ids["tag.tag1"])
	if err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)",
		ids["note.note2"], ids["tag.tag2"])
	if err != nil {
		return
	}

	return
}
