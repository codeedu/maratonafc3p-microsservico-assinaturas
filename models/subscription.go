package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
)

// Subscription is used by pop to map your subscriptions database table to your go code.
type Subscription struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	SubscriberID         uuid.UUID  `db:"subscriber_id"`
	Subscriber           Subscriber `json:"-" belongs_to:"subscriber" db:"-"`
	PlanID               uuid.UUID  `db:"plan_id"`
	Plan                 Plan       `json:"-" belongs_to:"plan" db:"-"`
	RemotePlanID         string     `db:"remote_plan_id"`
	RemoteSubscriptionID string     `db:"remote_subscription_id"`
	StartDate            time.Time  `db:"start_date"`
	ExpiresAt            time.Time  `db:"expires_at"`
	Status               string     `db:"status"`
	Payments             Payments   `json:"-" has_many:"payments" db:"-"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s Subscription) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Subscriptions is not required by pop and may be deleted
type Subscriptions []Subscription

// String is not required by pop and may be deleted
func (s Subscriptions) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Subscription) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Subscription) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Subscription) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
