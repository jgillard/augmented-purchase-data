package transactioncategories

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type jsonStatus struct {
	Status string `json:"status"`
}

func (c *Server) statusHandler(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	payload := marshallResponse(jsonStatus{"OK"})
	res.Write(payload)
}
