package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/cpgillem/csnotes"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	syntax := "Syntax: db setup|teardown"

	if len(os.Args) > 1 {
		db, err := sql.Open("mysql", "notes_app:notes_app@/notes_app")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if os.Args[1] == "setup" {
			fmt.Println("Setting up DB...")
			err = csnotes.SetUpDB(db)
			if err != nil {
				panic(err)
			}
		} else if os.Args[1] == "teardown" {
			fmt.Println("Tearing down DB...")
			err = csnotes.TearDownDB(db)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(syntax)
		}

		fmt.Println("Done.")
	} else {
		fmt.Println(syntax)
	}

}
