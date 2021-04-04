package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"student/models"
	"time"

	"github.com/gin-gonic/gin"
)

func HomeLink(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"data": "Ok"})
}

type Student struct {
	EnrollmentNumber string `json:"enrollmentNumber"`
	Name             string `json:"name"`
	Class            string `json:"class"`
	Subject          string `json:"subject"`
	Age              int    `json:"age"`
}

type CreateBookInput struct {
	Name    string `json:"name" binding:"required"`
	Class   string `json:"class" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Age     uint   `json:"age" binding:"required"`
}

func CreateStudent(c *gin.Context) {
	var student []Student
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are missing"})
		return
	}

	enrollmentNumber := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(enrollmentNumber)

	if err := models.Session.Query("INSERT INTO students(name, class, subject, isdeleted, enrollmentnumber, age) VALUES(?, ?, ?, ?, ?, ?)",
		input.Name, input.Class, input.Subject, false, strconv.Itoa(int(enrollmentNumber)), input.Age).Exec(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := models.Session.Query("INSERT INTO students_by_isDeleted(name, class, subject, isdeleted, enrollmentnumber, age) VALUES(?, ?, ?, ?, ?,?)",
		input.Name, input.Class, input.Subject, false, strconv.Itoa(int(enrollmentNumber)), input.Age).Exec(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	m := map[string]interface{}{}
	query := "SELECT * FROM students WHERE enrollmentnumber=?"
	iterable := models.Session.Query(query, strconv.Itoa(int(enrollmentNumber))).Iter()
	found := false
	for iterable.MapScan(m) {
		found = true
		student = append(student, Student{
			EnrollmentNumber: m["enrollmentnumber"].(string),
			Name:             m["name"].(string),
			Subject:          m["subject"].(string),
			Class:            m["class"].(string),
			Age:              m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	if found {
		c.JSON(http.StatusOK, gin.H{"data": student})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil})
	}

}

func GetAllStudents(c *gin.Context) {

	var students []Student
	m := map[string]interface{}{}
	query := "SELECT * FROM students_by_isdeleted where isDeleted=false;"
	iterable := models.Session.Query(query).Iter()
	found := false

	for iterable.MapScan(m) {
		found = true
		students = append(students, Student{
			EnrollmentNumber: m["enrollmentnumber"].(string),
			Name:             m["name"].(string),
			Subject:          m["subject"].(string),
			Class:            m["class"].(string),
			Age:              m["age"].(int),
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
	query := "SELECT * FROM students WHERE enrollmentNumber=? AND isDeleted=false"
	iterable := models.Session.Query(query, c.Param("id")).Iter()
	found := false
	for iterable.MapScan(m) {
		found = true
		student = append(student, Student{
			EnrollmentNumber: m["enrollmentnumber"].(string),
			Name:             m["name"].(string),
			Subject:          m["subject"].(string),
			Class:            m["class"].(string),
			Age:              m["age"].(int),
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

	var student []Student

	m := map[string]interface{}{}
	query := "SELECT * FROM students WHERE enrollmentNumber=? AND isDeleted=false"
	iterable := models.Session.Query(query, c.Param("id")).Iter()
	found := false
	for iterable.MapScan(m) {
		found = true
		student = append(student, Student{
			EnrollmentNumber: m["enrollmentnumber"].(string),
			Name:             m["name"].(string),
			Subject:          m["subject"].(string),
			Class:            m["class"].(string),
			Age:              m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	if found {
		if err := models.Session.Query("INSERT INTO deleted_students(name, class, subject, enrollmentnumber, age) VALUES(?, ?, ?, ?, ? )",
			student[0].Name, student[0].Class, student[0].Subject, student[0].EnrollmentNumber, student[0].Age).Exec(); err != nil {
			fmt.Printf("FAILED AT INSERT 1")
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"data": "Student not found"})
	}

	deleteQuery := "DELETE FROM students WHERE enrollmentNumber=?"
	if err := models.Session.Query(deleteQuery, c.Param("id")).Exec(); err != nil {
		fmt.Printf("FAILED AT DELETE 1")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	deleteQuery2 := "DELETE FROM Students_by_isDeleted WHERE isDeleted=false and enrollmentNumber=?"
	if err := models.Session.Query(deleteQuery2, c.Param("id")).Exec(); err != nil {
		fmt.Printf("FAILED AT DELETE 2")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"data": "Student deleted successfully"})
}
