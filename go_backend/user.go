package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func addNewUser(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var objMap map[string]interface{}
	err := json.Unmarshal(reqBody, &objMap)
	if err != nil {
		return
	}
	objMap["password"], _ = HashPassword(objMap["password"].(string))
	newReqBody, err := json.Marshal(objMap)

	var user User
	json.Unmarshal(newReqBody, &user)

	//write to db
	sqlStatement := `
		INSERT INTO users (email, first_name, last_name, password)
		VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, objMap["email"], objMap["first_name"], objMap["last_name"], objMap["password"])
	if err != nil {
		panic(err)
	}
	db.Close()

	json.NewEncoder(w).Encode(user)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var objMap map[string]interface{}
	err := json.Unmarshal(reqBody, &objMap)
	if err != nil {
		return
	}
	//Check name exists and get password
	userExists, password := checkUser(objMap["email"].(string))
	if !userExists {
		w.WriteHeader(404)
		return
	}
	//Compare the two
	if !CheckPasswordHash(objMap["password"].(string), password) {
		w.WriteHeader(401)
		return
	}
	db.Close()

	token, tokErr := GenerateJWT(objMap["email"].(string))
	if tokErr != nil {
		w.WriteHeader(401)
		return
	}
	w.Write([]byte(token))
	//Create a JWT if success
}

func getSprintList(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["id"]

	sqlStatementBoards := `select * from board where user_id = $1`
	results, err := db.Query(sqlStatementBoards, id)
	if err != nil {
		panic(err.Error())
	}

	var boardList []Board
	for results.Next() {
		var b Board
		// for each row, scan the result into our tag composite object
		err = results.Scan(&b.BoardID, &b.StartDate, &b.EndDate, &b.IsActive, &b.User)

		if err != nil {
			panic(err.Error())
		}
		b.User = ""
		boardList = append(boardList, b)
	}

	objMap := make(map[string]interface{})
	objMap["sprintBoards"] = boardList

	json.Marshal(objMap)
	json.NewEncoder(w).Encode(objMap)

	db.Close()
}
