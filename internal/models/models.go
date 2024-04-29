package models

type Car struct {
	RegNum string `json:"regNum" query:"regNum"`
	Mark   string `json:"mark" query:"mark"`
	Model  string `json:"model" query:"model"`
	Owner  People `json:"people ,omitempty"`
	Year   int    `json:"year ,omitempty"  query:"year"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic ,omitempty"`
}
