package transactioncategories

import "encoding/json"

type jsonID struct {
	ID string `json:"id"`
}

type jsonName struct {
	Name string `json:"name"`
}

type jsonTitle struct {
	Title string `json:"title"`
}

type jsonStatus struct {
	Status string `json:"status"`
}

type jsonError struct {
	Title string `json:"title"`
}

type jsonErrors struct {
	Errors []jsonError `json:"errors"`
}

const (
	contentTypeKey = "Content-Type"
	statusDeleted  = "deleted"
)

func jsonIsValid(body []byte) bool {
	var js struct{}
	return json.Unmarshal(body, &js) == nil
}
