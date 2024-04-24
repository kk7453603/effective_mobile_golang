package models

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Owner  uint64 `json:"owner"`
}

type People struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
