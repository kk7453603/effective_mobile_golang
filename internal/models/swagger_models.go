package models

type Mock_Response struct {
}

type Response_Car struct {
	RegNum string `json:"regNum,omitempty"`
	Mark   string `json:"mark,omitempty"`
	Model  string `json:"model,omitempty"`
	Year   int    `json:"year ,omitempty"`
}

type Response_User struct {
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Response_Error struct {
	Message string `json:"message"`
}

type Response_OK struct {
	Status string `json:"status"`
}

type Request_Add_Cars struct {
	RegNums []string `json:"regNums"`
}
