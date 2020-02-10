package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "HabitTracker"
)

var mySigningKey = []byte("secret")

func connectToDataBase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	//fmt.Print(reflect.TypeOf(db))
	return db
}

func generateUUID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)
	return string(out)
}

//HashPassword
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//GenerateJWT is a helper function that creates JSON web tokens
func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

//Check existance of
func checkUser(username string) (bool, string) {
	db := connectToDataBase()
	var passwordHash string
	sqlStatement := `select password from users where email = $1`
	err := db.QueryRow(sqlStatement, username).Scan(&passwordHash)
	db.Close()
	if err != nil {
		return false, ""
	}

	return true, passwordHash
}

func checkHabit(habitID string) bool {
	db := connectToDataBase()
	sqlStatement := `select * from habit where id = $1`
	_, err := db.Query(sqlStatement, habitID)
	db.Close()
	if err != nil {
		return false
	}
	return true
}

func checkBoard(boardID string) bool {
	db := connectToDataBase()
	sqlStatement := `select * from board where id = $1`
	_, err := db.Query(sqlStatement, boardID)
	db.Close()
	if err != nil {
		return false
	}
	return true
}
