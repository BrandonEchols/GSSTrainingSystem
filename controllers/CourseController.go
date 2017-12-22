package controllers

import (
	"GSSTrainingSystem/models"
	"GSSTrainingSystem/services"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
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

const GET_PAGE_TEMPLATE string = "%s?activity=%d"

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

	_, activity, act_err := this.course_service.GetCourseAndActivity(course_name, activity_num)
	if act_err != nil {
		fmt.Println("Error returned from GetCourseActivity. Err: ", act_err.Error())
		this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
		return
	}
	_ = activity.GetType()
	//TODO check the posted request against the activity returned

	activity_num++ //increment to next page

	fmt.Println("Redirecting to next page...")

	r.Method = "GET"
	new_url := fmt.Sprintf(GET_PAGE_TEMPLATE, strings.Split(r.URL.Path, "?")[0], activity_num)
	http.Redirect(w, r, new_url, 302)
}

type LayoutTemplateData struct {
	Title        string
	ActivityHead template.HTML
	ActivityBody template.HTML
}

type VideoBodyTemplateData struct {
	Url template.URL
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
		new_url := fmt.Sprintf(GET_PAGE_TEMPLATE, r.URL.Path, 0)
		http.Redirect(w, r, new_url, 302)
		return
	}

	activity_num, err := strconv.Atoi(activity_string)
	if err != nil {
		this.WriteErrorMessageWithStatus(w, 400, "Invalid 'activity' query parameter")
		return
	}
	course, activity, act_err := this.course_service.GetCourseAndActivity(course_name, activity_num)
	if act_err != nil {
		fmt.Println("Error returned from GetCourseActivity. Err: ", act_err.Error())
		this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
		return
	}

	layout_data := LayoutTemplateData{
		Title: course.Title,
	}

	switch asserted_data := activity.(type) {
	case *models.StaticActivity:
		bodyPath := "assets" + asserted_data.HtmlPath
		headPath := strings.TrimSuffix(bodyPath, ".html") + "-head.html"

		body_bytes, err := ioutil.ReadFile(bodyPath)
		if err != nil {
			fmt.Println("Could not read filepath: ", bodyPath)
			this.WriteErrorMessageWithStatus(w, 500, "Could not read file "+bodyPath)
			return
		}

		head_bytes, err := ioutil.ReadFile(headPath)
		if err != nil {
			fmt.Println("Could not read filepath: ", headPath)
			this.WriteErrorMessageWithStatus(w, 500, "Could not read file "+headPath)
			return
		}

		layout_data.ActivityBody = template.HTML(body_bytes)
		layout_data.ActivityHead = template.HTML(head_bytes)

	case *models.VideoActivity:
		body_path := "templates/video-activity.html"
		head_path := "templates/video-activity-head.html"
		video_body_data := VideoBodyTemplateData{
			Url: template.URL(asserted_data.VideoUrl),
		}

		buf := new(bytes.Buffer)
		var body_t *template.Template
		body_t = template.Must(body_t.ParseFiles(body_path))
		err = body_t.Execute(buf, video_body_data)
		if err != nil {
			fmt.Printf("Unable to parse %s file. Err: %s", body_path, err.Error())
			this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
			return
		}
		layout_data.ActivityBody = template.HTML(buf.String())

		buf = new(bytes.Buffer)
		var head_t *template.Template
		head_t = template.Must(head_t.ParseFiles(head_path))
		err = head_t.Execute(buf, nil)
		if err != nil {
			fmt.Printf("Unable to parse %s file. Err: %s", head_path, err.Error())
			this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
			return
		}
		layout_data.ActivityHead = template.HTML(buf.String())

	case *models.MultipleChoiceActivity:

		//TODO Once Multi-Choice has been implemented, implement this VV

		/*
			head_path := "templates/video-activity-head.html"
			body_path := "templates/video-activity.html"

			buf := new(bytes.Buffer)
			var body_t *template.Template
			body_t = template.Must(body_t.ParseFiles(body_path))
			err = body_t.Execute(buf, video_body_data)
			if err != nil {
				fmt.Printf("Unable to parse %s file. Err: %s", body_path, err.Error())
				this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
				return
			}
			layout_data.ActivityBody = template.HTML(buf.String())

			buf = new(bytes.Buffer)
			var head_t *template.Template
			head_t = template.Must(head_t.ParseFiles(head_path))
			err = head_t.Execute(buf, nil)
			if err != nil {
				fmt.Printf("Unable to parse %s file. Err: %s", head_path, err.Error())
				this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
				return
			}
			layout_data.ActivityHead = template.HTML(buf.String())
		*/
	default:
		fmt.Println("Unknown activity returned from Course Service")
		this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
		return
	}

	var t *template.Template
	t = template.Must(t.ParseFiles("templates" + course.Layout))
	err = t.Execute(w, layout_data)
	if err != nil {
		fmt.Println("Unable to parse layout file. Err: ", err.Error())
		this.WriteErrorMessageWithStatus(w, 500, "internal_server_error")
		return
	}
}
