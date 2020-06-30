package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
)

// Subscriber is used by pop to map your subscribers database table to your go code.
type Subscriber struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	Name           string        `json:"name" db:"name"`
	Email          string        `json:"email" db:"email"`
	DocumentNumber string        `json:"document_number" db:"document_number"`
	Street         string        `json:"street" db:"street"`
	StreetNumber   string        `json:"street_number" db:"street_number"`
	Complementary  string        `json:"complementary" db:"complementary"`
	Neighborhood   string        `json:"neighborhood" db:"neighborhood"`
	Zipcode        string        `json:"zipcode" db:"zipcode"`
	DDD            string        `json:"ddd" db:"ddd"`
	Number         string        `json:"number" db:"number"`
	Subscriptions  Subscriptions `has_many:"subscriptions" db:"-"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s Subscriber) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Subscribers is not required by pop and may be deleted
type Subscribers []Subscriber

// String is not required by pop and may be deleted
func (s Subscribers) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Subscriber) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Subscriber) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Subscriber) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
