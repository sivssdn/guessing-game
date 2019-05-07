package handler

import(
  "database/sql"

	"github.com/kataras/iris"
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
  session_id,_ := session(ctx).GetInt("session_id")

  if session_id >= 0{
    ctx.ViewData("session_id", session_id)
    ctx.View("player1.html")
  }else{
    ctx.Redirect("/", iris.StatusTemporaryRedirect)
  }
}
