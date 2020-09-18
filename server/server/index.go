package main

import (
	"fmt"
	"log"
	"net/http"

	"./controllers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {
	fmt.Println("Started")
	controllers.Configer()
	controllers.Database()
	r := mux.NewRouter()
	r.HandleFunc("/SignIn", controllers.SignIn).Methods("POST")
	r.HandleFunc("/SignUp", controllers.SignUp).Methods("POST")
	r.HandleFunc("/Movies", controllers.Movies).Methods("POST")
	r.HandleFunc("/locations", controllers.Locations).Methods("GET")
	r.HandleFunc("/theatres", controllers.Theatres).Methods("POST")
	r.HandleFunc("/getTickets", controllers.GetTickets).Methods("POST")
	r.HandleFunc("/pay", controllers.IndexHandler).Methods("GET")
	r.HandleFunc("/callback", controllers.CallBackHandler).Methods("POST")

	// Solves Cross Origin Access Issue
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":3001",
	}

	log.Fatal(srv.ListenAndServe())
}
