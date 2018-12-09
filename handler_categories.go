package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/xid"
)

type CategoryStore interface {
	GetCategoryList() CategoryList
}

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *CategoryServer) CategoryGetHandler(res http.ResponseWriter, req *http.Request) {
	cats := &CategoryList{
		Categories: []Category{
			{ID: "a1b2", Name: "foo"},
		},
	}

	bytePayload, err := json.Marshal(cats)

	if err != nil {
		log.Fatal(err)
	}

	res.Write(bytePayload)
	res.WriteHeader(http.StatusOK)
}

func (c *CategoryServer) CategoryPostHandler(res http.ResponseWriter, req *http.Request) {

	guid := xid.New().String()

	res.WriteHeader(http.StatusCreated)
	payload := fmt.Sprintf(`{"id":"%s","name":"accommodation"}`, guid)
	res.Write([]byte(payload))
}
