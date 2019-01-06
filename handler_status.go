package transactioncategories

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const statusBodyString = "OK"

func (c *Server) statusHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(res, statusBodyString)
}
