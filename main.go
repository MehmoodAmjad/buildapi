package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for course (goes in file)
type Course struct {
	CourseID    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fake DB
var courses []Course

// middleware, helper(goes in file)
func (c *Course) IsEmpty() bool {
	//return c.CourseID == "" && c.CourseName == ""
	return c.CourseName == ""
}

func main() {
	fmt.Println("Welcome to building an api inn GOLANG")
	r := mux.NewRouter()

	// seeding
	courses = append(courses, Course{CourseID: "2", CourseName: "GOLANG", CoursePrice: 299, Author: &Author{Fullname: "Mehmood Amjad", Website: "securiti.go"}})
	courses = append(courses, Course{CourseID: "4", CourseName: "Docker", CoursePrice: 399, Author: &Author{Fullname: "Mehmood Amjad", Website: "foundri.go"}})

	// routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{courseid}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{courseid}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{courseid}", deleteOneCourse).Methods("DELETE")

	// listen to port
	log.Fatal(http.ListenAndServe(":4000", r))

}

// conntrollers (goes in seperate files)

// serve home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to API</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all Courses")
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get One Course")
	w.Header().Set("Content=Type", "application/json")
	// grab id from request
	params := mux.Vars(r)
	// loop through courses and find matchingn id then return the reponse
	for _, course := range courses {
		if course.CourseID == params["courseid"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("no Course Founnd with given id")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create onne course")
	w.Header().Set("Content=Type", "application/json")

	// what if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	// what about data (being sent in form of {})
	var course Course

	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside JSON")
		return
	}

	// gennerate unnique id, string
	// append course into courses

	rand.Seed(time.Now().UnixNano())
	course.CourseID = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content=Type", "application/json")

	// first - grab id fro request
	params := mux.Vars(r)

	// loop through value to get id then remove then add with ID

	for index, course := range courses {
		if course.CourseID == params["courseid"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseID = params["courseID"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	// send a response whenn id not found
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one course")
	w.Header().Set("Content=Type", "application/json")

	params := mux.Vars(r)
	// loop, finnd id, remove(index,index+1)
	for index, course := range courses {
		if course.CourseID == params["courseid"] {
			courses = append(courses[:index], courses[index+1:]...)
			break
		}
	}
}
