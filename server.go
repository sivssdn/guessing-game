package main

import (

		"github.com/kataras/iris"
		"github.com/kataras/iris/middleware/recover"
		"github.com/r3labs/sse"

		handler "guessing-game/handler"
)

func main(){
 app := iris.New()
 app.Use(recover.New())

 //error handling
 app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context){
	 ctx.JSON(iris.Map{"error": "page not found"})
 })
 app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context){
	 ctx.JSON(iris.Map{"error": "internal server error"})
 })


 app.RegisterView(iris.HTML("./views", ".html").Reload(true))
 //for serving static css, js
 app.StaticWeb("/public", "./public")

 app.Handle("GET", "/", func(ctx iris.Context){
	 	ctx.View("index.html")
	});

 //for server sent events
 s := sse.New()

 app.Any("/player-1", handler.RenderPlayer1Page)
 app.Post("/player-1/login", func(ctx iris.Context){
	 	handler.Player1_login(ctx, s)
 })

 app.Any("/player-2", handler.RenderPlayer2Page)
 app.Post("/player-2/guess", handler.Guess)
 app.Post("/player-2/login",func(ctx iris.Context){
			handler.Player2_login(ctx, s)
 })

	app.Post("/player-1-question", func(ctx iris.Context){
		handler.PlayerStreamQuestions(ctx, s)
	})
	app.Post("/player-2-question", func(ctx iris.Context){
		handler.Player2Questions(ctx, s)
	})

	app.Any("/player-messages", iris.FromStd(s.HTTPHandler))

 app.Run(iris.Addr(":8180"))
}
