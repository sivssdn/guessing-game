package main

import (
		"github.com/kataras/iris"
		"github.com/kataras/iris/middleware/recover"
		"github.com/r3labs/sse"

		handler "guessing-game/handler"
		"time"
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

 app.Any("/player-1", handler.RenderPlayer1Page)
 app.Post("/player-1/login", handler.Player1_login)

 app.Any("/player-2", handler.RenderPlayer2Page)
 app.Post("/player-2/login", handler.Player2_login)
 app.Post("/player-2/guess", handler.Guess)


	s := sse.New()
	s.CreateStream("messages")
	app.Any("/player-question", iris.FromStd(s.HTTPHandler))

	go func() {

		for {
				time.Sleep(2 * time.Second)
				now := time.Now()
				s.Publish("messages", &sse.Event{
					Data: []byte(now.String()),
				})
		}
	}()

 app.Run(iris.Addr(":8180"))
}
