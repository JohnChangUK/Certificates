package model

import (
	"time"
)

const (
	Pending  = "PENDING"
	Complete = "COMPLETE"
	Declined = "DECLINED"
)

type Certificate struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	OwnerId   string    `json:"owner_id"`
	Year      int       `json:"year"`
	Note      string    `json:"note"`
	Transfer  *Transfer `json:"transfer"`
}

type Transfer struct {
	To     string `json:"email"`
	Status string `json:"status"`
}

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (transfer Transfer) TransferComplete(to string) *Transfer {
	return &Transfer{to, Complete}
}

func (transfer Transfer) TransferDeclined() *Transfer {
	return &Transfer{Status: Declined}
}
