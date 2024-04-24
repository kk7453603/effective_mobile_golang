package models

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Owner  People `json:"people"`
	Year   int    `json:"year"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}
