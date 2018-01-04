package services

import (
	"GSSTrainingSystem/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type CourseService struct{}

func GetCourseService() CourseService {
	return CourseService{}
}

func (this CourseService) GetCourseAndActivity(course_name string, activity_number int) (models.Course, models.IActivity, error) {
	activity_number_str := strconv.Itoa(activity_number)
	filePath := "courses/" + course_name + ".json"

	// Look for course_name as a file in the courses dir.
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		return models.Course{}, &models.Activity{}, errors.New("Could not open file " + filePath + " for reading. Err: " + err.Error())
	}

	// Read in that file
	data_bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return models.Course{}, &models.Activity{}, errors.New("Could not read in file " + filePath + ". Err: " + err.Error())
	}

	// Parse into Course Model
	course := models.Course{}
	err = json.Unmarshal(data_bytes, &course)
	if err != nil || len(course.Activities) < 1 {
		return models.Course{}, &models.Activity{}, errors.New("Could not parse json from file " + filePath + ". Err: " +
			err.Error())
	}

	if len(course.Activities) < activity_number {
		msg := fmt.Sprintf("Course %s only has %d activities, cannot get activity number %d",
			course_name,
			len(course.Activities),
			activity_number,
		)
		return models.Course{}, &models.Activity{}, errors.New(msg)
	}

	// Look for activity_number in that course
	activity_interface := course.Activities[activity_number]
	if activity_interface == nil {
		return models.Course{}, &models.Activity{}, errors.New("Activity number " + activity_number_str +
			" not found in course " + filePath)
	}

	// Parse that activity and determine it's type
	data, _ := activity_interface.(map[string]interface{})
	activity_type, ok := data["type"].(string)
	if !ok {
		return models.Course{}, &models.Activity{}, errors.New("Could not assert type of activity. Activity number " +
			activity_number_str + " in course " + filePath)
	}

	// Parse the activity as the correct type and return it.
	switch strings.ToUpper(activity_type) {
	case models.ACTIVITY_TYPE_STATIC:
		return_data := models.StaticActivity{}
		return_data.HtmlPath, ok = activity_interface.(map[string]interface{})["url"].(string)
		if !ok {
			return models.Course{}, &models.StaticActivity{}, errors.New("Unable to parse static activity." +
				" Activity number " + activity_number_str + " in course " + filePath)
		}
		return course, &return_data, nil
	case models.ACTIVITY_TYPE_VIDEO:
		return_data := models.VideoActivity{}
		return_data.VideoUrl, ok = activity_interface.(map[string]interface{})["url"].(string)
		if !ok {
			return models.Course{}, &models.StaticActivity{}, errors.New("Unable to parse video activity." +
				" Activity number " + activity_number_str + " in course " + filePath)
		}
		return course, &return_data, nil
	case models.ACTIVITY_TYPE_MULT_CHOICE:
		msg := "Unable to parse MultipleChoice activity." +
			" Activity number " + activity_number_str + " in course " + filePath + " bad field: "
		return_data := models.MultipleChoiceActivity{}
		activity, ok := activity_interface.(map[string]interface{})
		if !ok {
			return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "activity")
		}

		return_data.Question, ok = activity["question"].(string)
		if !ok {
			return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "question")
		}

		return_data.Answers = map[string]string{}
		answers, ok := activity["answers"].(map[string]interface{})
		if !ok {
			return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "answers")
		}
		for key, value := range answers { //Assert each answer is a string
			a, ok := value.(string)
			if !ok {
				return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "answers")
			}
			return_data.Answers[key] = a
		}

		correct_answers, ok := activity["correct"].([]interface{})
		if !ok {
			return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "correct")
		}
		for _, answer := range correct_answers { //Assert each answer is a string
			a, ok := answer.(string)
			if !ok {
				return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "answers")
			}
			return_data.CorrectAnswer = append(return_data.CorrectAnswer, a)
		}

		return_data.WrongText, ok = activity["wrongtext"].(string)
		if !ok {
			return models.Course{}, &models.MultipleChoiceActivity{}, errors.New(msg + "wrongtext")
		}
		return course, &return_data, nil
	default:
		msg := fmt.Sprintf("Unknown type of activity. Type %v Activity number %s in couse %s",
			data,
			activity_number_str,
			filePath,
		)
		return models.Course{}, &models.Activity{}, errors.New(msg)
	}
}
