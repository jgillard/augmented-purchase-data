package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/xid"
)

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CategoriesHandler(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		CategoryGetHandler(res, req)
	case http.MethodPost:
		CategoryPostHandler(res, req)
	}

}

func CategoryGetHandler(res http.ResponseWriter, req *http.Request) {
	cats := &CategoryList{
		Categories: []Category{
			{ID: "a1b2", Name: "foo"},
		},
	}

	bytePayload, err := json.Marshal(cats)

	if err != nil {
		log.Fatal(err)
	}

	res.Header().Set("content-type", jsonContentType)
	res.Write(bytePayload)
	res.WriteHeader(http.StatusOK)
}

func CategoryPostHandler(res http.ResponseWriter, req *http.Request) {

	guid := xid.New().String()

	res.Header().Set("content-type", jsonContentType)
	res.WriteHeader(http.StatusCreated)
	payload := fmt.Sprintf(`{"id":"%s","name":"accommodation"}`, guid)
	res.Write([]byte(payload))
}
