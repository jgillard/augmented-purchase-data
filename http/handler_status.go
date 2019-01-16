package httptransport

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	internal "github.com/jgillard/practising-go-tdd/internal"
)

func (c *Server) statusHandler(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	status := internal.GetStatus()
	payload := marshallResponse(jsonStatus{status})
	res.Write(payload)
}
