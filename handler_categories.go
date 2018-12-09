package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListCategoriesHandler(res http.ResponseWriter, req *http.Request) {
	cats := &CategoryList{
		Categories: []Category{
			{ID: "a1b2", Name: "foo"},
		},
	}

	bytePayload, err := json.Marshal(cats)

	if err != nil {
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(bytePayload)

}
