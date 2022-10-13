package model

type Response struct {
	Success bool        `json:"success"`
	Error   error       `json:"error"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
