package server

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ReturnError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	data, er := json.Marshal(&HttpError{
		Code: 500,
		Msg:  err.Error(),
	})

	if er != nil {
		w.Write([]byte("{code:500,msg:" + err.Error() + "}"))
	} else {
		w.Write(data)
	}
}
