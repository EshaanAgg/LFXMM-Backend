package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAllOrgs(c *gin.Context) {

	/*
	 * The following function creates a client instance for the database.
	 * GetAllParentOrgs function is called on this object instance.
	 *
	 * Defined at: ../database/handlers/client.go
	 * Returns:    Client object defined with function definition
	 */
	client := handlers.New()
	defer client.Close()

	/*
	 * The following function returns the fetched data.
	 * GetAllParentOrgs is defined at ../database/handlers/parentOrg.go
	 */
	c.IndentedJSON(http.StatusOK, client.GetAllParentOrgs())
}
