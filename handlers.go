package handlers

import (
	"io"
	"net/http"
)

const statusString = "OK"

func StatusHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, statusString)
}
