package response

import (
	"net/http"
)

const ErrHeader = "X-Error"

func WriteInternalError(w http.ResponseWriter, err error) {
	w.Header().Set(ErrHeader, err.Error())
	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
}
