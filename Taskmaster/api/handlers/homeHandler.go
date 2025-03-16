package handlers

import (
	"http"
)


func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	filePath := "./ReadMe.md"
	err := rnd.FileView(rw, http.StatusOk, filePath, "readme.md")
	CheckError(err)
}