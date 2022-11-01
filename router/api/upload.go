package api

import (
	"net/http"

	Conf "github.com/NotRoyadma/BDClient/config"
)

func UploadRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// check for the token header
	token := r.Header.Get("token")
	if token != Conf.Conf.ServiceToken {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Started Upload workers.`))
}
