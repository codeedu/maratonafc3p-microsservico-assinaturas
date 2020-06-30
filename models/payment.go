package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
)

// Payment is used by pop to map your payments database table to your go code.
type Payment struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	TransactionID  string       `json:"transaction_id" db:"transaction_id"`
	Gateway        string       `json:"gateway" db:"gateway"`
	PaymentType    string       `json:"payment_type" db:"payment_type"`

	CardBrand            string    `json:"card_brand" db:"card_brand"`
	CardLastDigits       string    `json:"card_last_digits" db:"card_last_digits"`
	BoletoURL            string    `json:"boleto_url" db:"boleto_url"`
	BoletoBarcode        string    `json:"boleto_barcode" db:"boleto_barcode"`
	BoletoExpirationDate string    `json:"boleto_expiration_date" db:"boleto_expiration_date"`

	Status         string       `json:"status" db:"status"`
	Total          int          `json:"total" db:"total"`
	Installments   int          `json:"installments" db:"installments"`
	Subscription   Subscription `belongs_to:"subscription" db:"-"`
	SubscriptionID uuid.UUID    `db:"subscription_id"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p Payment) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Payments is not required by pop and may be deleted
type Payments []Payment

// String is not required by pop and may be deleted
func (p Payments) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Payment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Payment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Payment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
