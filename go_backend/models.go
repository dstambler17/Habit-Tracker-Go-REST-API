package main

import "time"

//User - each user will be able to login
type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

//Board - the parent object, each habit must belong to a board
type Board struct {
	BoardID   string    `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	User      string    `json:"user"`
}

//Habit - the main object of the app
type Habit struct {
	HabitID     string `json:"id"`
	Title       string `json:"title"`
	Board       string `json:"board_id"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

//Note - a board can have many notes associated with it
type Note struct {
	NoteID string `json:"note_id"`
	Body   string `json:"content"`
	Board  string `json:"board_id"`
}
