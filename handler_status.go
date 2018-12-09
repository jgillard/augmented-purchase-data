package handlers

import (
	"io"
	"net/http"
)

const statusBodyString = "OK"

func (c *CategoryServer) statusHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, statusBodyString)
}
