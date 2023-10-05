package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProject(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	orgID := c.Param("id")
	projectID := c.Param("projectId")

	org := client.GetOrganizationByID(orgID)
	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this ID.",
		})
		return
	}

	// Using GetProjectById to get the project by project ID
	projects, err := client.GetProjectById(projectID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "No project was found with the provided ID.",
		})
		return
	}

	// Check if the organization ID associated with the project matches the supplied orgId
	if projects[0].OrganizationID != orgID {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "This project does not belong to the specified organization.",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"project": projects[0],
	})
}

func getProjectsByYear(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	id := c.Param("id")
	org := client.GetOrganizationByID(id)

	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this ID.",
		})
		return
	}

	yearParam := c.Query("year")

	// Check if yearParam is empty or not provided.
	if yearParam == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"projects": client.GetProjectsByOrganization(id),
		})
		return
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid year parameter. It should be of Integer type.",
		})
		return
	}

	projects := client.GetProjectsByYear(id, year)

	c.IndentedJSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

/*
 * The following function sends an API response with the number of projects by year
 * for a given organization (id)
 */
func getProjectCount(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	// Get id from HTTP request
	id := c.Param("id")

	// Check if an org exists with given id
	org := client.GetOrganizationByID(id)

	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this ID.",
		},
		)
		return
	}

	// Send out HTTP response. The function used below is defined under database/handlers
	c.IndentedJSON(http.StatusOK, client.GetCountOfProjectsByParentOrgID(id))
}
