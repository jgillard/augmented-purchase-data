package handlers

import (
	"io"
	"net/http"
)

const statusBodyString = "OK"

func StatusHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, statusBodyString)
}
