package api

import (
	"net/http"
	"fmt"
)

func HealthCheckHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "OK")
}