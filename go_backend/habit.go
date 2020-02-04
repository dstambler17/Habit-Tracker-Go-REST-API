package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func deleteHabit(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["id"]
	boardID := vars["board"]
	sqlStatement := `Delete from habit WHERE id = $1 and board_id = $2;`
	_, err := db.Exec(sqlStatement, id, boardID)
	if err != nil {
		panic(err)
	}
	db.Close()

}

func addHabit(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	boardID := vars["board"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	//fmt.Fprintf(w, "%+v", string(reqBody))

	//Add id and content into the object
	var objMap map[string]interface{}
	err := json.Unmarshal(reqBody, &objMap)
	if err != nil {
		return
	}
	objMap["id"] = strings.Replace(generateUUID(), "\n", "", 1)
	objMap["content"] = "{" + strings.Join([]string{"e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e"}, ",") + "}"
	newReqBody, err := json.Marshal(objMap)

	//convert json into the new object and add to model
	var habit Habit
	json.Unmarshal(newReqBody, &habit)

	sqlArray := "{" + strings.Join([]string{"e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e", "e"}, ",") + "}"
	//Add to db
	sqlStatement := `
		INSERT INTO habit (id, title, description, content, board_id)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, objMap["id"], objMap["title"], objMap["description"], sqlArray, boardID)
	if err != nil {
		panic(err)
	}
	db.Close()
	json.NewEncoder(w).Encode(habit)
}

func editHabit(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["id"]
	board := vars["board"]
	reqBody, _ := ioutil.ReadAll(r.Body)

	var objMap map[string]interface{}
	err := json.Unmarshal(reqBody, &objMap)
	if err != nil {
		return
	}

	sqlStatement := `update habit set title = $1, description = $2 where id = $3 and board_id = $4;`
	_, err = db.Exec(sqlStatement, objMap["title"], objMap["description"], id, board)
	if err != nil {
		panic(err)
	}

	db.Close()
}

func setStatus(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id, board, day, value := vars["id"], vars["board"], vars["day"], vars["val"]

	//get old string
	var oldArray string
	sqlStatement := `select content from habit where id = $1 and board_id = $2 limit 1`
	err := db.QueryRow(sqlStatement, id, board).Scan(&oldArray)
	if err != nil {
		panic(err)
	}

	//convert string to arr
	oldArray = strings.Replace(oldArray, "{", "", 1)
	oldArray = strings.Replace(oldArray, "}", "", 1)

	contentArr := strings.Split(oldArray, ",")
	i, err := strconv.Atoi(day)
	contentArr[i-1] = value
	newArrString := "{" + strings.Join(contentArr, ",") + "}"

	sqlStatementUpdate := `update habit set content = $1 where id = $2 and board_id = $3;`
	_, err = db.Exec(sqlStatementUpdate, newArrString, id, board)
	if err != nil {
		panic(err)
	}

	db.Close()
	w.Write([]byte(`{"new_content":` + newArrString + `}`))
	json.NewEncoder(w)
}

func getHabit(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	key := vars["id"]
	boardID := vars["board"]
	fmt.Println("Endpoint Hit: getHabit")

	//Query for the habit
	var habit Habit
	sqlStatement := `select * from habit where id = $1 and board_id = $2 limit 1`
	err := db.QueryRow(sqlStatement, key, boardID).Scan(&habit.HabitID, &habit.Title, &habit.Description, &habit.Content, &habit.Board)
	if err != nil {
		panic(err)
	}

	db.Close()
	json.NewEncoder(w).Encode(habit)
}
