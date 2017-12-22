package services

import (
	"GSSTrainingSystem/models"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type CourseService struct{}

func GetCourseService() CourseService {
	return CourseService{}
}

func (this CourseService) GetCourseActivity(course_name string, activity_number int) (models.IActivity, error) {
	activity_number_str := strconv.Itoa(activity_number)
	filePath := "courses/" + course_name + ".json"

	// Look for course_name as a file in the courses dir.
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		return &models.Activity{}, errors.New("Could not open file " + filePath + " for reading. Err: " + err.Error())
	}

	// Read in that file
	data_bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &models.Activity{}, errors.New("Could not read in file " + filePath + ". Err: " + err.Error())
	}

	// Parse into Course Model
	course := models.Course{}
	err = json.Unmarshal(data_bytes, &course)
	if err != nil {
		return &models.Activity{}, errors.New("Could not parse json from file " + filePath + ". Err: " +
			err.Error())
	}

	// Look for activity_number in that course
	activity_interface := course.Activities[activity_number]
	if activity_interface == nil {
		return &models.Activity{}, errors.New("Activity number " + activity_number_str +
			" not found in course " + filePath)
	}

	// Parse that activity and determine it's type
	data, _ := activity_interface.(map[string]interface{})
	activity_type, ok := data["type"].(string)
	if !ok {
		return &models.Activity{}, errors.New("Could not assert type of activity. Activity number " +
			activity_number_str + " in course " + filePath)
	}

	// Parse the activity as the correct type and return it.
	switch strings.ToUpper(activity_type) {
	case models.ACTIVITY_TYPE_STATIC:
		return_data := models.StaticActivity{}
		return_data.HtmlPath, ok = activity_interface.(map[string]interface{})["url"].(string)
		if !ok {
			return &models.StaticActivity{}, errors.New("Unable to parse static activity." +
				" Activity number " + activity_number_str + " in course " + filePath)
		}
		return &return_data, nil
	case models.ACTIVITY_TYPE_VIDEO:
		//TODO fix this VV
		return_data, ok := activity_interface.(models.VideoActivity)
		if !ok {
			return &models.VideoActivity{}, errors.New("Unable to parse video activity." +
				" Activity number " + activity_number_str + " in course " + filePath)
		}
		return &return_data, nil
	case models.ACTIVITY_TYPE_MULT_CHOICE:
		//TODO fix this VV
		return_data, ok := activity_interface.(models.MultipleChoiceActivity)
		if !ok {
			return &models.MultipleChoiceActivity{}, errors.New("Unable to parse MultipleChoice activity." +
				" Activity number " + activity_number_str + " in course " + filePath)
		}
		return &return_data, nil
	default:
		msg := fmt.Sprintf("Unknown type of activity. Type %v Activity number %s in couse %s",
			data,
			activity_number_str,
			filePath,
		)
		return &models.Activity{}, errors.New(msg)
	}
}
