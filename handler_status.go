package transactioncategories

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (c *server) statusHandler(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	payload := marshallResponse(jsonStatus{"OK"})
	res.Write(payload)
}
