package main

import (
	"fmt"
	"net/http"

	"github.com/rdnt/rdnt.dev/router"
)

func main() {
	fmt.Println(":)")

	//app := application.New()
	//_ = app

	r := router.New()

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
	fmt.Println("exit")
}
