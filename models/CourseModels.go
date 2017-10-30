package models

type Activity struct {
	CourseName string
	Number     int
	Type       string
}

type VideoActivity struct {
	Activity
	VideoUrl string
}

type MultipleChoiceActivity struct {
	Activity
	Question      string
	Answers       []string
	CorrectAnswer string
	BadResponse   string
	GoodResponse  string
}
