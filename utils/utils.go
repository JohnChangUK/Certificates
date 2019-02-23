package utils

import (
	. "github.com/JohnChangUK/verisart/model"
	"time"
)

func MockCertificates(certificates []Certificate) []Certificate {
	certificates = append(certificates, Certificate{Id: "1", Title: "First Certificate", CreatedAt: time.Now(),
		OwnerId: "John", Year: 2019, Note: "Art note",
		Transfer: &Transfer{}},
		Certificate{Id: "2", Title: "Second Certificate", CreatedAt: time.Now(),
			OwnerId: "Jim", Year: 2010, Note: "Painting note",
			Transfer: &Transfer{}})
	return certificates
}
