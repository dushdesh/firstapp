package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dushdesh/firstapp/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	// "github.com/nicholasjackson/env"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":3000", "Bind Address for the server")

func main() {
	l := log.New(os.Stdout, "first-app ", log.LstdFlags)

	// Intialize handlers
	user := handlers.NewUser(l)


	// Create new ServeMux and register route handlers
	sm := mux.NewRouter()

	// GET subrouter
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", user.GetAll)
	// sm.Handle("/users", u).Methods("Get")

	// POST subrouter
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", user.Create)
	postRouter.Use(user.MwValidateUser)

	// PUT subrouter
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", user.Update)
	putRouter.Use(user.MwValidateUser)

	// DELETE subrouter
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+", user.Delete)
	deleteRouter.Use(user.MwValidateUser)

	// Redoc document middleware to render swagger API file
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	rh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", rh)
	// Server the swagger file through the file server
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Create and configure a new Server
	s := http.Server {
		Addr: ":3000",
		Handler: sm,
		ErrorLog: l,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	// Start the server
	go func(){
		l.Println("Starting the server on port 3000")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting the server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Trap sigterm and gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block untill a signal is received
	sig := <-c
	log.Println("Got signal: ", sig)

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	s.Shutdown(ctx)
}

