package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sachinchaudhary003/MongoApi/router"
)

func main() {
	fmt.Println("Hello")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("listening port:8000")
}
