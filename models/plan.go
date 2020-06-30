package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
)

// Plan is used by pop to map your plans database table to your go code.
type Plan struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	Description   string        `json:"description" db:"description"`
	Price         float32       `json:"price" db:"price"`
	RemotePanID   string        `json:"remote_plan_id" db:"remote_plan_id"`
	Recurrence    string        `json:"recurrence" db:"recurrence"`
	Active        bool          `json:"active" db:"active"`
	Subscriptions Subscriptions `has_many:"subscriptions" db:"-"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p Plan) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Plans is not required by pop and may be deleted
type Plans []Plan

// String is not required by pop and may be deleted
func (p Plans) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Plan) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Plan) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Plan) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
