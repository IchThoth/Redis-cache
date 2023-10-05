package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ichthoth/Redis-cache/routes"
)

func Run() error {
	r := gin.Default()

	routes.UserRoutes(r)

	mount := r.Run()

	return mount
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
