package main

import (
	"github.com/leilei3167/cgin/framework"
	"log"
	"net/http"
)

func main() {
	server := http.Server{Handler: framework.NewCore(), Addr: ":8080"}
	log.Fatal(server.ListenAndServe())

}
