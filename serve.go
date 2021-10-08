//+build none

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Printf("Serving on http://localhost:3000/\n")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
