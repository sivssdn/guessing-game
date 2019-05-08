package model

type Word struct{
  ChoosenWord  string  `form:"word"`
}

type PlayerSession struct{
  SessionId   string  `form:"session_id"`
}

type PlayerInteraction struct{
  QuestionAnswer    string   `form:"question_ans"`
}
