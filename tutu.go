package main

import (
	"controllers"
	"flag"
	"strconv"

	"github.com/kataras/iris"
)

var (
	listenPort int
)

func init() {
	flag.IntVar(&listenPort, "p", 8001, "p=8001")
	flag.Parse()
}

func main() {
	iris.Get("/articles/:id", controllers.ArticlesGetHandler)
	iris.Post("/articles", controllers.ArticlesPostHandler)
	iris.Get("/attachs/:id", controllers.AttachsGetHandler)
	iris.Post("/attachs/:id", controllers.AttachsPostHandler)
	iris.Get("/tags/:id", controllers.TagsGetHandler)
	iris.Post("/tags/:id", controllers.TagsPostHandler)
	iris.Listen(":" + strconv.Itoa(listenPort))
}
