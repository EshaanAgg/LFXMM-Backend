package api

import (
	"eshaanagg/lfx/database/handlers"
	"fmt"
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

	// Using GetProjectByProjectId to get the project by project ID
	projects, err := client.GetProjectByProjectId(projectID)
	if err != nil {
		fmt.Print(projectID, orgID)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Project not found",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"id":   org.ID,
		"logo": org.Logo,
		"org":  org.Name,
		"project": gin.H{
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
		},
	})

}

func getProjectByYear(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	id := c.Param("id")
	org := client.GetOrganizationByID(id)

	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this id",
		},
		)
		return
	}

	yearParam := c.Query("year")
	Year, err := strconv.Atoi(yearParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
        return
    }

	projects := client.GetProjectsByYear(id,Year)

	c.IndentedJSON(http.StatusOK, gin.H{
		"projects":     projects,
	})
}

