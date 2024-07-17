// struct
package models

import (
	"gorm.io/gorm"
)

type Person struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	IsActive *bool  `json:"isactive" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email,min=3,max=32"`
}

type Register struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,usernameValidate"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	LineId   string `json:"lineid" validate:""`
	Tel      string `json:"tel" validate:"required,max=10"`
	Business string `json:"business" validate:"required"`
	Website  string `json:"website" validate:"required,min=2,max=30,websiteValidate"`
}

type Dogs struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
}

type DogsRes struct {
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Type  string `json:"type"`
}

type ResultData struct {
	Count      int       `json:"count"`
	Data       []DogsRes `json:"data"`
	Name       string    `json:"name"`
	SumRed     int       `json:"SumRed"`
	SumGreen   int       `json:"SumGreen"`
	SumPink    int       `json:"SumPink"`
	SumNoColor int       `json:"SumNoColor"`
}

type CompanyData struct {
	gorm.Model
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Position string `json:"position"`
	Salary   int    `json:"salary"`
}

type ProfileUser struct {
	gorm.Model
	Employee_id string `json:"employee_id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Birthday    string `json:"birthday"`
	Age         int    `json:"age"`
	Email       string `json:"email"`
	Tel         string `json:"tel"`
}

type AgeProfile struct {
	Employee string `json:"employee"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
	Gen      string `json:"gen"`
}

type ResultAgeProfile struct {
	Data    []AgeProfile `json:"data"`
	Count   int          `json:"count"`
	Name    string       `json:"name"`
	SumGenz int          `json:"countGenz"`
	SumGenY int          `json:"countGenY"`
	SumGenX int          `json:"countGenX"`
	SumBB   int          `json:"countBabyBoomer"`
	SumGI   int          `json:"countGI"`
}
