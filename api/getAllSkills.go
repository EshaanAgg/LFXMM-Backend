package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAllSkills(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	allSkills, err := client.GetAllSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, allSkills)
}
