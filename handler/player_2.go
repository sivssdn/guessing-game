package handler

import(
  "database/sql"

	"github.com/kataras/iris"
	_ "github.com/lib/pq"

  model "guessing-game/model"
)

func Player2_login(ctx iris.Context){

  session_model := &model.PlayerSession{}
  err := ctx.ReadForm(session_model)
  checkErr(err)

  dbCreds := DbCredentialsString()
	db, err := sql.Open("postgres", dbCreds)
	defer db.Close()
	checkErr(err)

  //check if the active session exists
  queryStmt, err := db.Prepare("SELECT session_id FROM public.sessions WHERE session_id=$1 AND game_status=1;")
  checkErr(err)
  var player1_session_id int
  err = queryStmt.QueryRow(session_model.SessionId).Scan(&player1_session_id)

  if err == sql.ErrNoRows {
    ctx.JSON(iris.Map{"status": "No Active Sessions Found"})
    return
  }
  checkErr(err) //error can be many

  session(ctx).Set("session_id", player1_session_id)

  ctx.Redirect("/player-2", iris.StatusTemporaryRedirect)

}

func RenderPlayer2Page(ctx iris.Context){
  session_id,_ := session(ctx).GetInt("session_id")

  if session_id >= 0{
    ctx.ViewData("session_id", session_id)
    ctx.View("player2.html")
  }else{
    ctx.Redirect("/", iris.StatusTemporaryRedirect)
  }
}

func Guess(ctx iris.Context){

  session_id,_ := session(ctx).GetInt("session_id")
  word := &model.Word{}
	err := ctx.ReadForm(word)
  checkErr(err)

  dbCreds := DbCredentialsString()
	db, err := sql.Open("postgres", dbCreds)
	defer db.Close()
	checkErr(err)

  queryStmt, err := db.Prepare("SELECT session_word FROM public.sessions WHERE session_id=$1 AND session_word=$2 AND game_status=1;")
  checkErr(err)
  var player1_session_word string
  err = queryStmt.QueryRow(session_id, word.ChoosenWord).Scan(&player1_session_word)

  if err == sql.ErrNoRows {
    ctx.JSON(iris.Map{"status": "Incorrect"})
  }
  checkErr(err) //error can be many

  if player1_session_word == word.ChoosenWord{
    //game over update status
    _, err = db.Exec("UPDATE public.sessions SET game_status=0 WHERE session_id=$1", session_id)
    checkErr(err)
    ctx.JSON(iris.Map{"status": "Correct"})
  }

}
