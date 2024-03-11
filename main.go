package main

import (
	"fmt"
	"net/http"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	myFigure := figure.NewFigure("Amin Nouri", "", true)
	myFigure.Print()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello World!")
}
