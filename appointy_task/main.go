package main

import (
	"Rest_api/operator"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/users", operator.Createuserendpoint)
	http.HandleFunc("/users/", operator.Userbyidendpoint)
	http.HandleFunc("/posts", operator.Createpostendpoint)
	http.HandleFunc("/posts/", operator.Getpostbyidendpoint)
	http.HandleFunc("/posts/users/", operator.Getuserspostbyidendpoint)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}