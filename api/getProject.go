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
			"message": "There is no organization with this id",
		})
		return
	}

	// Using GetProjectById to get the project by project ID
	projects, err := client.GetProjectById(projectID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Project not found",
		})
		return
	}

	// Check if the organization ID associated with the project matches the supplied orgId
	if projects[0].OrganizationID != orgID {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "This project does not belong to the specified organization",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"projectId":     projects[0].ProjectID,
		"lfxProjectUrl": projects[0].LFXProjectUrl,
		"name":          projects[0].Name,
		"industry":      projects[0].Industry,
		"description":   projects[0].Description,
		"repoLink":      projects[0].Repository,
		"websiteUrl":    projects[0].Website,
		"createdOn":     projects[0].CreatedOn,
		"amountRaised":  projects[0].AmountRaised,
		"skills":        projects[0].Skills,
		"parentOrg":     projects[0].OrganizationID,
	})
}

func getProjectsByYear(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	id := c.Param("id")
	org := client.GetOrganizationByID(id)

	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this id",
		})
		return
	}

	yearParam := c.Query("year")

	// Check if yearParam is empty or not provided.
	if yearParam == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Year parameter is required"})
		return
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}

	projects := client.GetProjectsByYear(id, year)

	c.IndentedJSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}
