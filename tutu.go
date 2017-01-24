package main

import (
	"controllers"

	"github.com/kataras/iris"
)

func main() {
	iris.Config.IsDevelopment = true // this will reload the templates on each reques
	iris.Get("/articles/:id", controllers.ArticlesGetHandler)
	iris.Post("/articles", controllers.ArticlesPostHandler)
	iris.Get("/attachs/:id", controllers.AttachsGetHandler)
	iris.Post("/attachs/:id", controllers.AttachsPostHandler)
	iris.Get("/tags/:id", controllers.TagsGetHandler)
	iris.Post("/tags/:id", controllers.TagsPostHandler)
	iris.Listen(":8080")
}
