package projectid

type hitResponse struct {
	ID string `json:"_id"`
}

type hitsResponse struct {
	Hits []hitResponse `json:"hits"`
}

type apiResponse struct {
	Hits hitsResponse `json:"hits"`
}
