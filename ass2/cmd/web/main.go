package main

import (
	"ass2"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/exchange", fotsopp.HandlerProjects)
	http.HandleFunc("/exchange/", fotsopp.HandlerWebhookId)
	http.HandleFunc("/exchange/latest", fotsopp.HandlerLatest)
	http.HandleFunc("/exchange/average", fotsopp.HandlerAverage)
	http.HandleFunc("/exchange/evaluationtrigger", fotsopp.HandlerEvaluation)
	http.ListenAndServe(":"+port, nil)
}
