package models

type Course struct {
	CourseName string        `json:"courseName"`
	Title      string        `json:"title"`
	Layout     string        `json:"layout"`
	Activities []interface{} `json:"activities"`
}

type IActivity interface {
	GetType() string
}

const ACTIVITY_TYPE_STATIC string = "STATIC"
const ACTIVITY_TYPE_VIDEO string = "VIDEO"
const ACTIVITY_TYPE_MULT_CHOICE string = "MULTIPLECHOICE"

type Activity struct {
	Type string `json:"type"`
}

type VideoActivity struct {
	Activity
	VideoUrl string `json:"url"`
}

type StaticActivity struct {
	Activity
	HtmlPath string `json:"url"`
}

type MultipleChoiceActivity struct {
	Activity
	Question      string            `json:"question"`
	Answers       map[string]string `json:"answers"`
	CorrectAnswer []string          `json:"correct"`
	WrongText     string            `json:"wrongtext"`
}

func (this Activity) GetType() string {
	return ""
}

func (this VideoActivity) GetType() string {
	return ACTIVITY_TYPE_VIDEO
}

func (this MultipleChoiceActivity) GetType() string {
	return ACTIVITY_TYPE_MULT_CHOICE
}

func (this StaticActivity) GetType() string {
	return ACTIVITY_TYPE_STATIC
}
