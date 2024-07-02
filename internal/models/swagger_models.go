package models

type Response_User struct {
	PassportNumber string `json:"passport_number" form:"passport_number"`
	Surname        string `json:"surname" form:"surname"`
	Name           string `json:"name" form:"name"`
	Patronymic     string `json:"patronymic" form:"patronymic"`
	Address        string `json:"address" form:"address"`
}

type Response_Error struct {
	Error string `json:"error"`
}

type Response_OK struct {
	Status string `json:"status"`
}
