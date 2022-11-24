package models

import "time"

type EchoDTO struct {
	Message string `json:"message"`
}

type EmailModel struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type ClientRequest struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	EmailPhone string    `json:"email_phone"`
	Request    string    `json:"request"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CRDeleteDTO struct {
	Id int `json:"id"`
}

type CheckNameEmailPhoneDTO struct {
	Name       string `json:"name"`
	EmailPhone string `json:"email_phone"`
}

type CheckResponseDTO struct {
	NameCheck  bool `json:"name_check"`
	EmailCheck bool `json:"email_check"`
	PhoneCheck bool `json:"phone_check"`
}
