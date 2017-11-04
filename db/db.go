package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/cpgillem/csnotes"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	syntax := "Syntax: db setup|teardown|regenerate|seed"
	valid := false

	if len(os.Args) > 1 {
		db, err := sql.Open("mysql", "notes_app:notes_app@/notes_app")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if os.Args[1] == "teardown" || os.Args[1] == "regenerate" {
			valid = true
			fmt.Println("Tearing down DB...")
			err = csnotes.TearDownDB(db)
			if err != nil {
				panic(err)
			}
		}

		if os.Args[1] == "setup" || os.Args[1] == "regenerate" {
			valid = true
			fmt.Println("Setting up DB...")
			err = csnotes.SetUpDB(db)
			if err != nil {
				panic(err)
			}
		}
		
		if os.Args[1] == "seed" {
			valid = true
			fmt.Println("Seeding DB...")
			_, err = csnotes.SeedDB(db)
			if err != nil {
				panic(err)
			}
		}
		
		if !valid {
			fmt.Println(syntax)
		} else {
			fmt.Println("Done.")
		}

	} else {
		fmt.Println(syntax)
	}

}
