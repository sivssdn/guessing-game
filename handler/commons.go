package handler

import(
  "fmt"
  "os"

  "github.com/joho/godotenv"
  "github.com/kataras/iris"
  "github.com/kataras/iris/sessions"
)

//session management
var (
    cookieNameForSessionID = "guessing_game_session_id_hijack?"
    sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true})
)
func session(ctx iris.Context){
    return sess.Start(ctx)
}

//for error checking
func checkErr(err error) {
        if err != nil {
						//do error logging here
						fmt.Println(err.Error())
						panic(err)
        }
}

//for database
func DbCredentialsString() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"))
}
