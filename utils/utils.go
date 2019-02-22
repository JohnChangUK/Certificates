package utils

import (
	. "github.com/JohnChangUK/verisart/models"
	"time"
)

func MockCertificates(certificates []Certificate) []Certificate {
	user := User{"userId1", "user1@gmail.com", "User1"}
	user2 := User{"userId2", "user2@gmail.com", "User2"}
	certificates = append(certificates, Certificate{"1", "First Certificate", time.Now(),
		"John", 2019, "Art note",
		&Transfer{user.Email, "Processing Transfer"}},
		Certificate{"2", "Second Certificate", time.Now(),
			"Jim", 2010, "Painting note",
			&Transfer{user2.Email, "Complete"}})
	return certificates
}
