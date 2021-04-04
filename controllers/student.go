package controllers

import (
	"fmt"
	"net/http"
	"student/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func HomeLink(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"data": "Ok"})
}

type Student struct {
	ID        gocql.UUID `json:"ID"`
	Firstname string     `json:"Firstname"`
	Lastname  string     `json:"Lastname"`
	Age       int        `json:"Age"`
	IsDeleted bool       `json:"isDeleted"`
}

type CreateBookInput struct {
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	Age       int    `json:"age" binding:"required"`
}

func CreateStudent(c *gin.Context) {
	var newStudent []Student
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are missing"})
		return
	}

	cqlUID := gocql.TimeUUID()
	if err := models.Session.Query("INSERT INTO students(id, firstname, lastname, age, isDeleted) VALUES(?, ?, ?, ?,?)",
		cqlUID, input.Firstname, input.Lastname, input.Age, false).Exec(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	m := map[string]interface{}{}
	query := "SELECT * FROM students WHERE id=?"
	iterable := models.Session.Query(query, cqlUID).Iter()
	found := false
	for iterable.MapScan(m) {
		found = true
		newStudent = append(newStudent, Student{
			ID:        m["id"].(gocql.UUID),
			Firstname: m["firstName"].(string),
			Lastname:  m["lastName"].(string),
			Age:       m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	if found {
		c.JSON(http.StatusOK, gin.H{"data": newStudent})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil})
	}

}

func GetAllStudents(c *gin.Context) {

	var students []Student
	m := map[string]interface{}{}
	query := "SELECT * FROM students WHERE isDeleted=false"
	iterable := models.Session.Query(query).Iter()
	found := false

	for iterable.MapScan(m) {
		found = true
		students = append(students, Student{
			ID:        m["id"].(gocql.UUID),
			Firstname: m["firstName"].(string),
			Lastname:  m["lastName"].(string),
			Age:       m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	if found {
		c.JSON(http.StatusOK, gin.H{"data": students})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil})
	}
}

func GetOneStudent(c *gin.Context) {
	var student []Student

	m := map[string]interface{}{}
	query := "SELECT * FROM students WHERE id=? AND isDeleted=false"
	iterable := models.Session.Query(query, c.Param("id")).Iter()
	found := false
	for iterable.MapScan(m) {
		found = true
		student = append(student, Student{
			ID:        m["id"].(gocql.UUID),
			Firstname: m["firstName"].(string),
			Lastname:  m["lastName"].(string),
			Age:       m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	if found {
		c.JSON(http.StatusOK, gin.H{"data": student})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil})
	}
}

func DeleteStudent(c *gin.Context) {
	query := "UPDATE STUDENTS SET isDeleted=true WHERE id=?"
	if err := models.Session.Query(query, c.Param("id")).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"data": "Student deleted successfully"})
}
