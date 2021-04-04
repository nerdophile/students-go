package main

import (
	"student/controllers"
	"student/models"

	"github.com/gin-gonic/gin"
)

/* Initial Database Queries

Primary Student Table
CREATE TABLE students (enrollmentNumber text, name text, class text, subject text, isDeleted boolean, PRIMARY KEY( enrollmentNumber, isDeleted));

For quering with isDeleted field
CREATE TABLE students_by_isDeleted (enrollmentNumber text, name text, class text, subject text, isDeleted boolean, PRIMARY KEY(isDeleted, enrollmentNumber));

Soft deleted students
CREATE TABLE deleted_students (enrollmentNumber text, name text, class text, subject text, PRIMARY KEY(enrollmentNumber));

*/

func main() {
	models.ConnectDatabase()
	router := gin.Default()
	router.GET("/", controllers.HomeLink)

	router.POST("/students", controllers.CreateStudent) // http://localhost:3000/students POST
	router.GET("/students", controllers.GetAllStudents) // http://localhost:3000/students GET
	// router.GET("/students/count", controllers.CountAllStudents)   // http://localhost:3000/students/count GET
	router.GET("/students/:id", controllers.GetOneStudent)   // http://localhost:3000/students/1 GET
	router.DELETE("/student/:id", controllers.DeleteStudent) // http://localhost:3000/students/1 DELETE
	// router.HandleFunc("/students", controllers.DeleteAllStudents).Methods("DELETE")     // http://localhost:3000/students/ DELETE
	// router.PATCH("/students/{id}", controllers.UpdateStudent) // http://localhost:3000/students/3 PATCH
	router.Run(":3000")
}
