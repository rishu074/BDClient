package api

import (
	"net/http"

	Conf "github.com/NotRoyadma/BDClient/config"
	Stats "github.com/NotRoyadma/BDClient/stats"
	Workers "github.com/NotRoyadma/BDClient/workers"
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

	stat, _ := Stats.GetStats()
	w.WriteHeader(http.StatusOK)
	if stat {
		w.Write([]byte(`Worker is already running.`))
		return
	}

	w.Write([]byte(`Started upload workers.`))
	Stats.SetStats(true)
	go Workers.StartUploadWorker()

}
