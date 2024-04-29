package models

type Car struct {
	RegNum string `json:"regNum" query:"regNum" form:"regNum"`
	Mark   string `json:"mark" query:"mark" form:"mark"`
	Model  string `json:"model" query:"model" form:"model"`
	Owner  People `json:"people ,omitempty" form:"people"`
	Year   int    `json:"year ,omitempty"  query:"year" form:"year"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic ,omitempty"`
}
