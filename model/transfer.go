package model

import . "github.com/JohnChangUK/Certificates/constants"

type Transfer struct {
	To     string `json:"email"`
	Status string `json:"status"`
}

func (transfer Transfer) TransferComplete(to string) *Transfer {
	return &Transfer{to, Complete}
}

func (transfer Transfer) TransferDeclined() *Transfer {
	return &Transfer{Status: Declined}
}
