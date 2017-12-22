package controllers

import (
	"GSSTrainingSystem/models"
	"GSSTrainingSystem/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

/*
	The CourseController Class handles requests to the "/courses/" family of endpoints.
*/
type CourseController struct {
	course_service services.CourseService
	controller
}

func GetCourseController(course_service services.CourseService) CourseController {
	return CourseController{
		course_service: course_service,
	}
}

/*
	POST -> /courses/<course_name>?activity=<activity_number>
	PostPage handles the requests that are meant to verify that the activity has been completed.
	If the data passed in the request meets the requirements, then the request is redirected
	to GET -> /courses/<course_name>?activity=<activity_number>++
	If the data passed in is invalid, a "400" error is returned with a message explaining the incorrect data.
*/
func (this CourseController) PostPage(w http.ResponseWriter, r *http.Request) {
	course_name := strings.TrimPrefix(r.URL.Path, "/courses/")
	if course_name == "" {
		this.WriteErrorMessageWithStatus(w, 400, "Unable to determine course name")
		return
	}

	activity_string := r.URL.Query().Get("activity")
	if activity_string == "" {
		this.WriteErrorMessageWithStatus(w, 400, "Missing 'activity' query parameter")
		return
	}

	activity_num, err := strconv.Atoi(activity_string)
	if err != nil {
		this.WriteErrorMessageWithStatus(w, 400, "Invalid 'activity' query parameter")
		return
	}

	activity, act_err := this.course_service.GetCourseActivity(course_name, activity_num)
	if act_err != nil {
		fmt.Println("Error returned from GetCourseActivity. Err: ", act_err.Error())
		this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
		return
	}
	_ = activity.GetType()
	//TODO check the posted request against the activity returned

	activity_num++ //increment to next page
	nextPageNum := strconv.Itoa(activity_num)

	fmt.Println("Redirecting to next page...")

	r.Method = "GET"
	const GET_PAGE_TEMPLATE string = "/courses/%s?activity=%s"
	new_url := fmt.Sprintf(GET_PAGE_TEMPLATE, strings.Split(r.URL.Path, "?")[0], nextPageNum)
	http.Redirect(w, r, new_url, 302)
}

/*
	GET -> /courses/<course_name>?activity=<activity_number>
	GetPage returns the standard page with the activity asked for in the URL path injected in it.
	It reads the 'courses' directory looking for the <course_name>.json file. If no file exists, a 404 Not Found is
	returned.
	If the file exists, it looks for the activity_number requested, if the activity does not exist, a 404 Not Found
	is returned.
	If the activity exists, it determines what type of activity to inject into the page and returns it.
*/
func (this CourseController) GetPage(w http.ResponseWriter, r *http.Request) {
	course_name := strings.TrimPrefix(r.URL.Path, "/courses/")
	if course_name == "" {
		this.WriteErrorMessageWithStatus(w, 400, "Unable to determine course name")
		return
	}

	activity_string := r.URL.Query().Get("activity")
	if activity_string == "" {
		this.WriteErrorMessageWithStatus(w, 400, "Missing 'activity' query parameter")
		return
	}

	activity_num, err := strconv.Atoi(activity_string)
	if err != nil {
		this.WriteErrorMessageWithStatus(w, 400, "Invalid 'activity' query parameter")
		return
	}
	activity, act_err := this.course_service.GetCourseActivity(course_name, activity_num)
	if act_err != nil {
		fmt.Println("Error returned from GetCourseActivity. Err: ", act_err.Error())
		this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
		return
	}

	switch asserted_data := activity.(type) {
	case *models.StaticActivity:
		http.Redirect(w, r, "/assets/"+asserted_data.HtmlPath, 302)
		return
	case *models.VideoActivity:
		//TODO do something
		return
	case *models.MultipleChoiceActivity:
		//TODO do something
		return
	}

	fmt.Println("Unknown activity returned from Course Service")
	this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
}
