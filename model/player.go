package model

type Word struct{
  ChoosenWord  string  `form:"word"`
}

type PlayerSession struct{
  SessionId   string  `form:"session_id"`
}
