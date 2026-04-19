package models

type Card struct {
	ID           int    `json:"id"`
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	CorrectCount int    `json:"correctCount"`
	WrongCount   int    `json:"wrongCount"`
}
