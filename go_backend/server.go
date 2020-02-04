package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/habit/{board}", addHabit).Methods("POST")
	myRouter.HandleFunc("/habit/{board}/{id}", getHabit).Methods("GET")
	myRouter.HandleFunc("/habit/{board}/{id}", deleteHabit).Methods("DELETE")
	myRouter.HandleFunc("/habit/{board}/{id}/edit", editHabit).Methods("PUT")
	myRouter.HandleFunc("/habit/{board}/{id}/status/{day}/{val}", setStatus).Methods("PUT")
	myRouter.HandleFunc("/board", addBoard).Methods("POST")
	myRouter.HandleFunc("/board/{id}", getBoardDetails).Methods("GET")
	myRouter.HandleFunc("/board/{id}", deleteBoard).Methods("DELETE")
	myRouter.HandleFunc("/board/{id}", changeActive).Methods("PUT")
	myRouter.HandleFunc("/note/{board}", editNote).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

	fmt.Print("hi")
	//fmt.Print(reflect.TypeOf(generateUUID()))*/
	handleRequests()
}
