package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userid"`
	TaskName  string    `json:"taskname"`
	Content   string    `json:"content"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type User struct {
	ID             int    `json:"id"`
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}
