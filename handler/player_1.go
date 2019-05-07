package handler

import(
  "database/sql"

	// "encoding/json"
	"github.com/kataras/iris"
  "github.com/kataras/iris/sessions"
	_ "github.com/lib/pq"

  model "guessing-game/model"
)

func Player1_login(ctx iris.Context){

	word := &model.Word{}
	err := ctx.ReadForm(word)
  checkErr(err)

  dbCreds := DbCredentialsString()
	db, err := sql.Open("postgres", dbCreds)
	defer db.Close()
	checkErr(err)

  var session_id int
	err = db.QueryRow("INSERT INTO public.sessions(session_word) VALUES ($1) RETURNING session_id;", word.ChoosenWord).Scan(&session_id)
	checkErr(err)

  session(ctx).Set("session_id", session_id)
  ctx.Redirect("/player-1", iris.StatusTemporaryRedirect)

}

func RenderPlayer1Page(ctx iris.Context){
  print(session(ctx).Get("session_id"))
  print("----------------------------")

  ctx.View("player1.html")
}
