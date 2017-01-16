package main

import (
	"net/http"

	r "./router"
)

func main() {
	router := r.GetRouter()
	http.ListenAndServe(":8081", router)
}
