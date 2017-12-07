# Notes App

A simplified Google-keep-like notes app. Submitted as a semester project for CS 5550 (Networks) at Western Michigan University.

# Prerequisites

- [Go 1.9](https://golang.org/doc/install)
- [MySQL](https://mysql.com)

# Setup Guide

1. Install the prerequisites:

   ```bash
   $ go get -u github.com/gorilla/mux
   $ go get -u github.com/dgriijalva/jwt-go
   $ go get -u github.com/auth0/go-jwt-middleware
   $ go get -u github.com/urfave/negroni
   $ go get -u github.com/go-sql-driver/mysql
   ```
   
1. Generate RSA keys for the app directory:
   ```bash
   scripts/keygen.sh app/keys
   ```
   
1. Change to the directory of the project and bulid it:

   ```bash
   $ cd notes-app-cs5550
   $ go build
   ```
   
1. Create a database user in MySQL with the username `notes_app` and the password `notes_app`.

1. Set the database up and seed it:

   ```bash
   $ db/db setup
   $ db/db seed
   ```
   
1. Run the app from the `app` directory:

   ```bash
   $ cd app
   $ ./app 8080
   ```
   
1. Access the page from `http://localhost:8080/`. Log in with the username`nonadmin` and password `password`.

# Dev Environment Notes

- Create database and user
  
  database: `notes_app`
  
  user: `notes_app`
- Run migration scripts in db/

# References

- [Practical Persistence in Go: Organising Database Access](http://www.alexedwards.net/blog/organising-database-access)
- [Building Web Apps with Go](https://codegangsta.gitbooks.io/building-web-apps-with-go/content/)
- [Interface With Your Database in Go Tests](https://robots.thoughtbot.com/interface-with-your-database-in-go)
- [Go unit test setup and teardown](https://blog.karenuorteva.fi/go-unit-test-setup-and-teardown-db1601a796f2)
- https://stackoverflow.com/questions/25218903/how-are-people-managing-authentication-in-go
- https://www.youtube.com/watch?v=dgJFeqeXVKw
- https://gist.github.com/thealexcons/4ecc09d50e6b9b3ff4e2408e910beb22
- https://crackstation.net/hashing-security.htm
- https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
- http://blog.restcase.com/rest-api-error-codes-101/
- https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
- https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#4xx_Client_errors
- https://gowebexamples.com/templates/
- https://elithrar.github.io/article/approximating-html-template-inheritance/
- https://scotch.io/tutorials/create-a-single-page-app-with-go-echo-and-vue
