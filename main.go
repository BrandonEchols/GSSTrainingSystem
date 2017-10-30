package main

import (
	"GSSTrainingSystem/controllers"
	"GSSTrainingSystem/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	course_service := services.GetCourseService()
	course_service.ReadInCourses()

	course_controller := controllers.GetCourseController(course_service)

	router.Methods("POST").PathPrefix("/courses/").HandlerFunc(course_controller.PostPage)
	router.Methods("GET").PathPrefix("/courses/").HandlerFunc(course_controller.GetPage)
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", router)
}
