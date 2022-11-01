package router

import (
	"net/http"
	"strings"

	"github.com/NotRoyadma/BDClient/logger"
	"github.com/NotRoyadma/BDClient/router/api"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Write the HTTP logs
	logger.WriteAutoHTTPLogs(w, r)

	//Handle Different Paths
	if strings.Contains(r.URL.Path, "/upload") {
		api.UploadRequestHandler(w, r)
		return
	}

	//404 page
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
