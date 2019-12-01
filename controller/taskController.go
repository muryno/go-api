package controller

import (
	"go-api/model"
	u "go-api/utils"
	"net/http"
)

var CreateTask = func(w http.ResponseWriter, r *http.Request) {
	var P model.Task

	if r.Body == nil {
		u.Responds(w, u.Message(false, "Please send a request body"))
		return
	}

	resp := P.AddTask()
	u.Responds(w, resp)

}
