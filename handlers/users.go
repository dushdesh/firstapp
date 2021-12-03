// Package Manage Users API
//
// Documentation for Users API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dushdesh/firstapp/data"
	"github.com/gorilla/mux"
)

// List of users returned in the response
// swagger:response usersResponse
type usersResponse struct {
	// All the users in the system
	// in: body
	Body []data.User
}

type User struct {
	l *log.Logger
}

func NewUser(l *log.Logger) *User {
	return &User{l}
}

// func (u *User) ServeHTTP(rw http.ResponseWriter, r *http.Request){

// 	// Handle the GET request
// 	if r.Method == http.MethodGet {
// 		u.getAll(rw, r)
// 		return
// 	}

// 	// Handle POST request
// 	if r.Method == http.MethodPost {
// 		u.create(rw, r)
// 		return
// 	}

// 	// Handle PUT request
// 	if r.Method == http.MethodPut {
// 		reg := regexp.MustCompile(`/([0-9]+)`)
// 		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
// 		if len(g) != 1 {
// 			http.Error(rw, "Invalid URI no id", http.StatusBadRequest)
// 			return
// 		}
// 		if len(g[0]) != 2 {
// 			http.Error(rw, "Invalid URI more than 2 capture groups", http.StatusBadRequest)
// 			return
// 		}
// 		idString := g[0][1]
// 		id, err := strconv.Atoi(idString)
// 		if err != nil {
// 			http.Error(rw, "Invalid URI could not convert string ID to ID", http.StatusBadRequest)
// 			return
// 		}
// 		u.update(id, rw, r)
// 	}

// 	// catch all unsupported methods
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

// swagger:route GET /users users listUsers
// Returns a list of users
// responses:
// 	200: usersResponse
func (u *User) GetAll(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET request for Users")
	users := data.GetUsers()

	err := users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall JSON", http.StatusInternalServerError)
	}
}

func (u *User) Create(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle CREATE request")

	// Get the user from the context set in the middleware
	user := r.Context().Value(KeyUser{}).(*data.User)

	// Add the User to the User list
	ul := data.AddUser(user)

	// Marshall the user list to the response
	err := ul.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall appended user list", http.StatusInternalServerError)
		return
	}
}

func (u *User) Update(rw http.ResponseWriter, r *http.Request) {

	u.l.Println("Handle UPDATE request")

	// Fetch the value of id from the url
	id, err := getId(r)
	if err != nil {
		http.Error(rw, "Bad ID value provided", http.StatusBadRequest)
		return
	}

	// Get the user from the context set in the middleware
	nu := r.Context().Value(KeyUser{}).(*data.User)

	// Make sure the url ID and the object ID is same in case the request JSON has a different ID than url params
	nu.Id = id

	// Replace the old user with new user
	err = data.UpdateUser(id, nu)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
		return
	}
}

func (u *User) Delete(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle DELETE request")

	// Fetch the value of id from the url
	id, err := getId(r)
	if err != nil {
		http.Error(rw, "Bad ID value provided", http.StatusBadRequest)
		return
	}

	err = data.DeleteUser(id)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to delete user", http.StatusInternalServerError)
		return
	}
}

func getId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	return id, err
}

type KeyUser struct{}

func (u *User) MwValidateUser(next http.Handler) http.Handler {

	u.l.Println("Middelware: User Validation")

	return http.HandlerFunc( func(rw http.ResponseWriter, r *http.Request) {

		nu := &data.User{}

		// Unmarshall request data JSON to user
		err := nu.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user")
			http.Error(
				rw,
				"Unable to unmarshall user json",
				http.StatusInternalServerError,
			)
			return
		}

		// Validate the user
		err = nu.Validate()
		if err != nil {
			u.l.Println("[ERROR] user validation failed")
			http.Error(
				rw,
				fmt.Sprintf("User data is invalid: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// Add the unmarshalled user to the context to be passed on forward
		ctx := context.WithValue(r.Context(), KeyUser{}, nu)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(rw, r)
	})
}
