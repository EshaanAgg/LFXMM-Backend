package project

import (
	"eshaanagg/lfx/database/handlers"
	"fmt"
)

var getIdStmt = `SELECT id FROM parentOrgs WHERE name = $1`
var updateStmt = `UPDATE projects SET organizationId = $1 WHERE organizationId = $2`
var deleteStmt = `DELETE FROM parentOrgs WHERE id = $1`

func RemoveCNCF() {
	client := handlers.New()
	defer client.Close()

	allOrgs := client.GetAllOrgNames()
	for _, org := range allOrgs {
		l := len(org)
		if l > 7 && org[0:7] == "CNCF - " {
			Rename(org[7:l], org)
		}
	}
}

func Rename(newName string, oldName string) {
	client := handlers.New()
	defer client.Close()

	client.QueryRow(`UPDATE parentOrgs SET name = $1 WHERE name = $2`, newName, oldName)
}

func Merge() {
	client := handlers.New()
	defer client.Close()

	toMerge := make([][]string, 10)
	for i := 0; i < 10; i++ {
		toMerge[i] = make([]string, 10)
	}

	for _, mergeList := range toMerge {
		mainOrg := mergeList[0]
		var mainOrgId string

		err := client.QueryRow(getIdStmt, mainOrg).Scan(&mainOrgId)
		if err != nil {
			fmt.Println(err)
		}

		for i := 1; i < len(mergeList); i++ {
			org := mergeList[i]
			var orgId string

			err := client.QueryRow(getIdStmt, org).Scan(&orgId)
			if err != nil {
				fmt.Println(err)
			}

			client.QueryRow(updateStmt, mainOrgId, orgId)
			client.QueryRow(deleteStmt, orgId)
		}
	}
}
