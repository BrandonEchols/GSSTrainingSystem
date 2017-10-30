package controllers

import (
	"GSSTrainingSystem/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	activity_string := r.URL.Query().Get("activity")
	if activity_string != "" {
		this.WriteErrorMessageWithStatus(w, 400, "Missing 'activity' query parameter")
		return
	}

	activity_num, err := strconv.Atoi(activity_string)
	if err != nil {
		this.WriteErrorMessageWithStatus(w, 400, "Invalid 'activity' query parameter")
		return
	}

	_, _ = this.course_service.GetActivity(activity_num)
	//TODO check the posted request against the activity returned

	activity_num++ //increment to next page
	nextPageNum := strconv.Itoa(activity_num)

	fmt.Println("Redirecting to next page...")

	http.Redirect(w, r, "/courses/"+nextPageNum, 302)
}

/*
	GET -> /courses/<course_name>/<activity_number>
	GetPage returns the standard page with the activity asked for in the URL path injected in it.
	It reads the 'courses' directory looking for the <course_name>.json file. If no file exists, a 404 Not Found is
	returned.
	If the file exists, it looks for the activity_number requested, if the activity does not exist, a 404 Not Found
	is returned.
	If the activity exists, it determines what type of activity to inject into the page and returns it.
*/
func (this CourseController) GetPage(w http.ResponseWriter, r *http.Request) {
	activity_string := r.URL.Query().Get("activity")
	if activity_string != "" {
		this.WriteErrorMessageWithStatus(w, 400, "Missing 'activity' query parameter")
		return
	}

	activity_num, err := strconv.Atoi(activity_string)
	if err != nil {
		this.WriteErrorMessageWithStatus(w, 400, "Invalid 'activity' query parameter")
		return
	}
	_, _ = this.course_service.GetActivity(activity_num)
	//TODO Inject the activity into the appropriate template to return

	marshaled_data, _ := json.Marshal(activity_num)
	w.Write(marshaled_data)
}
