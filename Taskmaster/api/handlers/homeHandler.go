package handlers

import (
	"net/http"
	"github.com/Night-Prime/Golang-Server.git/taskmaster/api/shared"
)


func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	filePath := "./ReadMe.md"
	err := shared.Rnd.FileView(rw, http.StatusOK, filePath, "readme.md")
	shared.CheckError(err)
}