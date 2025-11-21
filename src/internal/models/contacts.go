package models

import (
	"time"
)

type Contacts struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Active    bool       `json:"active"`
	Linkedin  string     `json:"linkedin"`
	Github    string     `json:"github"`
	Credly    string     `json:"credly"`
}

type ContactsDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	LinkedIn string `json:"linkedin"`
	Github   string `json:"github"`
	Credly   string `json:"credly"`
}

func ToContactsDTO(contact *Contacts) *ContactsDTO {
	return &ContactsDTO{
		Name:     contact.Name,
		Email:    contact.Email,
		LinkedIn: contact.Linkedin,
		Github:   contact.Github,
		Credly:   contact.Credly,
	}
}
