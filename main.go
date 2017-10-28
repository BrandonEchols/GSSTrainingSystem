package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	router := mux.NewRouter()

	router.Methods("POST").PathPrefix("/courses/").HandlerFunc(postPage)
	router.Methods("GET").PathPrefix("/courses/").HandlerFunc(getPage)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", router)
}

func postPage(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/courses/"))
	if err != nil {
		fmt.Errorf("Could not parse %s to integer. Err: %v", strings.TrimPrefix(r.URL.Path, "/courses/"), err)
		w.WriteHeader(500)
		return
	}

	pageNum++ //increment to next page
	nextPageNum := strconv.Itoa(pageNum)

	fmt.Println("Redirecting to next page...")

	http.Redirect(w, r, "/courses/"+nextPageNum, 302)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageNum := strings.TrimPrefix(r.URL.Path, "/courses/")

	marshaled_data, _ := json.Marshal(pageNum)
	w.Write(marshaled_data)
}
