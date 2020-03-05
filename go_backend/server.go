package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
	myRouter.Handle("/habit/{board}", isAuthorized(addHabit)).Methods("POST")
	myRouter.Handle("/habit/{board}/{id}", isAuthorized(getHabit)).Methods("GET")
	myRouter.Handle("/habit/{board}/{id}", isAuthorized(deleteHabit)).Methods("DELETE")
	myRouter.Handle("/habit/{board}/{id}/edit", isAuthorized(editHabit)).Methods("PUT")
	myRouter.Handle("/habit/{board}/{id}/status/{day}/{val}", isAuthorized(setStatus)).Methods("PUT")
	myRouter.Handle("/board", isAuthorized(addBoard)).Methods("POST")
	myRouter.Handle("/board/{id}", isAuthorized(getBoardDetails)).Methods("GET")
	myRouter.Handle("/board/{id}", isAuthorized(deleteBoard)).Methods("DELETE")
	myRouter.Handle("/board/{id}", isAuthorized(changeActive)).Methods("PUT")
	myRouter.Handle("/note/{board}", isAuthorized(editNote)).Methods("PUT")
	myRouter.HandleFunc("/user/signup", addNewUser).Methods("POST")
	myRouter.HandleFunc("/user/login", logIn).Methods("PUT")
	myRouter.Handle("/user/getBoards/{id}", isAuthorized(getSprintList)).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error has occured")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func main() {
	//fmt.Print(reflect.TypeOf(generateUUID()))*/
	handleRequests()
}
