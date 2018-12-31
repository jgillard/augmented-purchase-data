package transactioncategories

import (
	"io"
	"net/http"
)

const statusBodyString = "OK"

func (c *Server) statusHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, statusBodyString)
}
