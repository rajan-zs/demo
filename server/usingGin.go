package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type employee struct {
	Name    string `json:"name"`
	Address string `json:"Address"`
}

func main() {
	r := gin.Default()
	r.GET("/getEmp", getEmp)
	r.GET("/postEmp", postEmp)
	r.Run()
}

func postEmp(c *gin.Context) {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		name := c.PostForm("name")
		add := c.PostForm("address")

		c.JSON(200, gin.H{
			"status":  "posted to employee details",
			"message": "whoo"})
		fmt.Println(name, add)
	})

}

func getEmp(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"name": "rajan",
	})
}
