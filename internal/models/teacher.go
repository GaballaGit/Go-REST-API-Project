package models

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Class     string `json:"class"`
	Subject   string `json:"subject"`
}
