package services

import "GSSTrainingSystem/models"

type CourseService struct{}

func GetCourseService() CourseService {
	return CourseService{}
}

func (this CourseService) GetActivity(activity_num int) (models.Activity, error) {
	//TODO Do some stuff and things.

	return models.Activity{}, nil
}

func (this CourseService) ReadInCourses() {
	//TODO Read in the *.json files in the courses directory for use.
}
