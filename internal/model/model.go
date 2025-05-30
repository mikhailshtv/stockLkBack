package model

type Success struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TokenSuccess struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Error struct {
	Error string `json:"error"`
}
