package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func getBoardDetails(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["id"]

	//Select board object
	var board Board
	sqlStatement := `select * from board where id = $1 limit 1`
	err := db.QueryRow(sqlStatement, id).Scan(&board.BoardID, &board.StartDate, &board.EndDate, &board.IsActive, &board.User)
	if err != nil {
		panic(err)
	}

	//Select Note object
	var note Note
	sqlStatementNote := `select * from note where board_id = $1 limit 1`
	err = db.QueryRow(sqlStatementNote, id).Scan(&note.NoteID, &note.Body, &note.Board)
	if err != nil {
		panic(err)
	}

	//Get the list of habits
	sqlStatementHabits := `select * from habit where board_id = $1`
	results, habitErr := db.Query(sqlStatementHabits, id)
	if habitErr != nil {
		panic(err.Error())
	}
	var habitList []Habit
	for results.Next() {
		var h Habit
		// for each row, scan the result into our tag composite object
		err = results.Scan(&h.HabitID, &h.Title, &h.Description, &h.Content, &h.Board)

		if err != nil {
			panic(err.Error())
		}
		h.Board = ""
		habitList = append(habitList, h)
	}

	//Add all to dict and return
	objMap := make(map[string]interface{})
	objMap["habits"] = habitList
	objMap["note"] = note
	objMap["id"] = board.BoardID
	objMap["start_date"] = board.StartDate
	objMap["end_date"] = board.EndDate
	objMap["is_active"] = board.IsActive
	objMap["user"] = board.User

	fmt.Println(objMap)

	json.Marshal(objMap)
	json.NewEncoder(w).Encode(objMap)

	db.Close()
}

func addBoard(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var objMap map[string]interface{}
	err := json.Unmarshal(reqBody, &objMap)
	if err != nil {
		return
	}

	objMap["id"] = strings.Replace(generateUUID(), "\n", "", 1)
	objMap["is_active"] = true

	//sprint to start at date of creation and will end in two weeks
	currentTime := time.Now()
	objMap["start_date"] = currentTime
	objMap["end_date"] = currentTime.AddDate(0, 0, 14)
	fmt.Println(objMap["start_date"])
	fmt.Println(objMap["end_date"])

	newReqBody, err := json.Marshal(objMap)

	var board Board
	json.Unmarshal(newReqBody, &board)

	//insert board
	sqlStatemenBoard := `
		INSERT INTO board (id, start_date, end_date, is_active, user_id)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatemenBoard, objMap["id"], objMap["start_date"], objMap["end_date"], objMap["is_active"], objMap["user_id"])
	if err != nil {
		panic(err)
	}

	//also create a new note for the board
	sqlStatementNote := `
	INSERT INTO note (id, body, board_id)
	VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatementNote, strings.Replace(generateUUID(), "\n", "", 1), "", objMap["id"])
	if err != nil {
		panic(err)
	}

	db.Close()
	json.NewEncoder(w).Encode(board)
}

func deleteBoard(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["id"]
	sqlStatement := `Delete from board WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	db.Close()
}

func editNote(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["board"]
	reqBody, _ := ioutil.ReadAll(r.Body)

	var objMap map[string]interface{}
	err := json.Unmarshal(reqBody, &objMap)
	if err != nil {
		return
	}

	//update note with the new body
	sqlStatement := `update note set body = $1 where board_id = $2;`
	_, err = db.Exec(sqlStatement, objMap["body"], id)
	if err != nil {
		panic(err)
	}

	db.Close()
}

func changeActive(w http.ResponseWriter, r *http.Request) {
	db := connectToDataBase()
	vars := mux.Vars(r)
	id := vars["id"]

	var isActive bool
	sqlStatement := `select is_active from board where id = $1 limit 1`
	err := db.QueryRow(sqlStatement, id).Scan(&isActive)
	if err != nil {
		panic(err)
	}

	sqlStatementUpdate := `update board set is_active = $1 where id = $2;`
	_, err = db.Exec(sqlStatementUpdate, !isActive, id)
	if err != nil {
		panic(err)
	}

	db.Close()
	w.Write([]byte(`{"new_active":` + strconv.FormatBool(!isActive) + `}`))
	json.NewEncoder(w)
}
