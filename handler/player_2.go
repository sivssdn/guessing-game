package handler

import(
  "database/sql"
  "strconv"

	"github.com/kataras/iris"
  "github.com/r3labs/sse"
	_ "github.com/lib/pq"

  model "guessing-game/model"
)

func Player2_login(ctx iris.Context, s *sse.Server){

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
  //for server sent events, question answers
  sse_id := strconv.Itoa(player1_session_id)+"message"
  s.Publish(sse_id, &sse.Event{
					Data: []byte("Connected for notifications - Player 2"),
	})

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
  }else{
      checkErr(err) //error can be many
  }


  if player1_session_word == word.ChoosenWord{
    //game over update status
    _, err = db.Exec("UPDATE public.sessions SET game_status=0 WHERE session_id=$1", session_id)
    checkErr(err)
    ctx.JSON(iris.Map{"status": "Correct"})
  }

}

func Player2Questions(ctx iris.Context, s *sse.Server){
  //checking if the limit of 20 question is over
  dbCreds := DbCredentialsString()
	db, err := sql.Open("postgres", dbCreds)
	defer db.Close()
	checkErr(err)
  session_id,_ := session(ctx).GetInt("session_id")
  //update questions count
  _, err = db.Exec("UPDATE public.sessions SET questions_count=questions_count+1 WHERE session_id=$1", session_id)
  checkErr(err)
  //check total questions questions_count
  queryStmt, err := db.Prepare("SELECT questions_count FROM public.sessions WHERE session_id=$1;")
  checkErr(err)
  var player_questions_count int
  err = queryStmt.QueryRow(session_id).Scan(&player_questions_count)
  checkErr(err)
  if player_questions_count > 20{
      //game over
      _, err = db.Exec("UPDATE public.sessions SET game_status=0 WHERE session_id=$1", session_id)
      checkErr(err)
      sse_id := strconv.Itoa(session_id)+"message"
      s.Publish(sse_id, &sse.Event{
              Data: []byte("Game Over!"),
      })
      
  }else{
      //display question
      PlayerStreamQuestions(ctx, s)
  }
}
