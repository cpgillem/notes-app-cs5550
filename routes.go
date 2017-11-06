package csnotes

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// GetIndex handles requests for the main page of the site.
func GetIndex(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	}
}

// GetLogin should take the user's credentials and create a
// JSON web token if the authentication was successful.
func PostLogin(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		
		// Validate the data.
		if len(username) == 0 {
			http.Error(w, "No username.", http.StatusInternalServerError)
			return
		}

		if len(password) == 0 {
			http.Error(w, "No password.", http.StatusInternalServerError)
			return
		}

		// Validate the credentials against the database.
		user, err := ValidateUser(username, password, context.DB)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could not authenticate user.", http.StatusForbidden)
			return
		}

		// Create a token.
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims {
			"iss": "admin",
			"exp": time.Now().Add(time.Minute * 20).Unix(),
			// TODO: Possibly use a generated struct that contains all necessary data
			"CustomUserInfo": struct {
				ID int64
			} {user.ID},
		})

		tokenString, err := token.SignedString(context.SignKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		response := struct {
			Token string `json:"token"`
		} {tokenString}

		// Turn the token into a json string.
		json, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func CreateRouter(context *Context) *mux.Router {
	router := mux.NewRouter()

	// Create middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return context.VerifyKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Public Routes
	router.HandleFunc("/", GetIndex(context))
	router.HandleFunc("/login", PostLogin(context)).Methods("POST")
	//router.HandleFunc("/logout", GetLogout(context)).Methods("GET")

	router.HandleFunc("/user", PostUser(context)).Methods("POST")
	router.HandleFunc("/user/{id}", GetUser(context)).Methods("GET")
	router.HandleFunc("/user/{id}", PutUser(context)).Methods("PUT")

	// Authenticated Routes
	router.Handle("/note/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(GetNote(context)),
	))

	return router
}
