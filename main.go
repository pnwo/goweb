package main

import (
	"net/http"
	"goweb/models"	
	"strconv"
	
	"github.com/gin-gonic/gin"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}



func getPersons(c *gin.Context) {

	persons, err := models.GetPersons(10)
	checkErr(err)
	if persons == nil {
		c.JSON(http.StatusOK, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": persons})
	}
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")

	person, err := models.GetPersonById(id)
	checkErr(err)

	if person.FirstName == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return 
	} else {
		c.JSON(http.StatusOK, gin.H{"data": person})
	}


	//c.JSON(http.StatusOK, gin.H{"message": "getPersonById" + id + "Called"})
}

func addPerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddPerson(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})	
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}


	// c.JSON(http.StatusOK, gin.H{"message": "addPerson Called"})
}

func updatePerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})	
		return
	}
	personId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdatePerson(json, personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})	
	}

	// c.JSON(http.StatusOK, gin.H{"message": "updatePerson Called"})
}

func deletePerson(c *gin.Context) {

	personId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeletePerson(personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	

	// id := c.Param("id")
	// c.JSON(http.StatusOK, gin.H{"message": "deletePerson" + id + "Called"})
}

func options(c *gin.Context) {
	
	ourOptions := "HTTP/1.1 200 OK\n" + 
		"Allow: GET, POST, PUT, DELETE, OPTIONS\n" + 
		"Access-Control-Allow-Origin: http://localhost:8080\n" + 
		"Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	c.String(200, ourOptions)
}
	// c.JSON(http.StatusOK, gin.H{"message": "options Called"})

func main() {

	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("person", getPersons)
		v1.GET("person/:id", getPersonById)
		v1.POST("person", addPerson)
		v1.PUT("person/:id", updatePerson)
		v1.DELETE("person/:id", deletePerson)
		v1.OPTIONS("person", options)
	}
	r.Run()

}