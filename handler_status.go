package transactioncategories

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const statusBodyJSON = `{"status":"OK"}`

func (c *Server) statusHandler(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	io.WriteString(res, statusBodyJSON)
}
