package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"subscription_service/models"
	"time"
)

// This is the struct responsible to aggregate all entities and services in order to process a new subscription
type PaymentService struct {
	Subscriber    models.Subscriber
	Subscription  models.Subscription
	Payment       models.Payment
	Connection    *pop.Connection
	PaymentReturn PaymentReturn
	ProcessData   ProcessData
	RabbitMQ      *RabbitMQ
}

// The PaymentReturn is the struct with the exact format which is received after a payment request is made
type PaymentReturn struct {
	ID                   string `json:"transaction_id"`
	Provider             *Gateway
	ProcessType          string `json:"object"`
	RemoteSubscriptionID int    `json:"id"`
	Status               string `json:"status"`
	CurrentTransaction   struct {
		RemoteTransactionID  int    `json:"id"`
		Amount               int    `json:"amount"`
		Installments         int    `json:"installments"`
		BoletoURL            string `json:"boleto_url"`
		BoletoBarcode        string `json:"boleto_barcode"`
		BoletoExpirationDate string `json:"boleto_expiration_date"`
	} `json:"current_transaction"`
	PaymentMethod      string    `json:"payment_method"`
	CardBrand          string    `json:"card_brand"`
	RemotePlanID       int       `json:"remote_plan_id"`
	PostbackURL        string    `json:"postback_url"`
	CardLastDigits     string    `json:"card_last_digits"`
	SoftDescriptor     string    `json:"soft_descriptor"`
	CurrentPeriodStart string    `json:"current_period_start"`
	CurrentPeriodSEnd  string    `json:"current_period_end"`
	RefuseReason       string    `json:"refuse_reason"`
	CreatedAt          time.Time `json:"date_created"`
	UpdatedAt          time.Time `json:"date_created"`
}

// ProcessData is responsible to bind the information sent via subscription
type ProcessData struct {
	Key            string    `json:"key" db:"-"`
	PlanID         uuid.UUID `json:"plan_id" db:"plan_id"`
	RemotePlanID   string    `json:"remote_plan_id" db:"remote_plan_id"`
	Name           string    `json:"name" db:"name"`
	Email          string    `json:"email" db:"email"`
	PaymentMethod  string    `json:"payment_method" db:"payment_method"`
	DocumentNumber string    `json:"document_number" db:"document_number"`
	CardHash       string    `json:"card_hash" db:"card_hash"`
	Street         string    `json:"street" db:"street"`
	StreetNumber   string    `json:"street_number" db:"street_number"`
	Complementary  string    `json:"complementary" db:"complementary"`
	Neighborhood   string    `json:"neighborhood" db:"neighborhood"`
	State          string    `json:"state" db:"state"`
	Zipcode        string    `json:"zipcode" db:"zipcode"`
	DDD            string    `json:"ddd" db:"ddd"`
	PhoneNumber    string    `json:"number" db:"number"`
}

// In order to make a payment request is necessary to inform the basic information about the gateway which is going to
// to process the request
type Gateway struct {
	ID            string `json:"-"`
	Name          string `json:"name"`
	ApiKey        string `json:"api_key"`
	EncryptionKey string `json:"encryption_key"`
}

// This is the required format to start a payment request
type TransactionSubscriptionRequest struct {
	SecretKey             string                `json:"secret_key"`
	Gateway               Gateway               `json:"gateway"`
	TransactionProviderID string                `json:"-"`
	APIKey                string                `json:"api_key"`
	RemotePlanID          int                   `json:"plan_id"`
	PaymentMethod         string                `json:"payment_method"`
	CardHash              string                `json:"card_hash"`
	SoftDescriptor        string                `json:"soft_descriptor"`
	PostbackURL           string                `json:"postback_url"`
	Customer              *CustomerSubscription `json:"customer"`
}

// This struct is responsible for handling customer information for transactions used for subscriptions only
type CustomerSubscription struct {
	ID             string `json:"-"`
	CustomerName   string `json:"name"`
	CustomerEmail  string `json:"email"`
	DocumentNumber string `json:"document_number"`
}

// Creates an empty PaymentService
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// Process the the subscription by doing:
// 1) Create the remote subscription
// 2) Create the subscription locally bu registering the customer data as well as the payment return information
// 3) Send to a queue the subscriber information
func (p *PaymentService) Process(data ProcessData) error {

	url := os.Getenv("PAYMENT_SUBSCRIPTION_ENDPOINT")
	p.ProcessData = data

	plan := models.Plan{}
	if err := p.Connection.Find(&plan, p.ProcessData.PlanID); err != nil {
		return err
	}
	p.ProcessData.RemotePlanID = plan.RemotePanID

	rPlanID, _ := strconv.Atoi(p.ProcessData.RemotePlanID)

	SubscriptionRequest := TransactionSubscriptionRequest{
		SecretKey: os.Getenv("PAYMENT_SECRET_KEY"),
		Gateway: Gateway{
			Name:   os.Getenv("GATEWAY"),
		},
		APIKey:         os.Getenv("GATEWAY_APIKEY"),
		RemotePlanID:   rPlanID,
		PaymentMethod:  p.ProcessData.PaymentMethod,
		CardHash:       p.ProcessData.CardHash,
		SoftDescriptor: "codeshop",
		PostbackURL:    "",
		Customer: &CustomerSubscription{
			CustomerName:   p.ProcessData.Name,
			CustomerEmail:  p.ProcessData.Email,
			DocumentNumber: p.ProcessData.DocumentNumber,
		},
	}

	jsonData, err := json.Marshal(SubscriptionRequest)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &p.PaymentReturn)

	if p.PaymentReturn.Status == "Declined" {
		log.Println("Transaction declined")
		return errors.New("Transaction declined")
	}
	err = p.insertData()

	if err != nil {
		log.Println("Error inserting data:", err)
		return err
	}

	if p.PaymentReturn.PaymentMethod == "credit_card" {
		subscriberJson, _ := json.Marshal(p.Subscriber)
		p.RabbitMQ.Notify(string(subscriberJson), "application/json", os.Getenv("RABBITMQ_NOTIFICATION_EX"), os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"))
	}

	return nil
}

func (p *PaymentService) insertData() error {

	subscriberId, _ := uuid.NewV4()
	subscriptionId, _ := uuid.NewV4()
	paymentId, _ := uuid.NewV4()

	startPeriod, _ := time.Parse(time.RFC3339, p.PaymentReturn.CurrentPeriodStart)
	endPeriod, _ := time.Parse(time.RFC3339, p.PaymentReturn.CurrentPeriodSEnd)
	// Subscription
	p.Subscription.ID = subscriptionId
	p.Subscription.SubscriberID = subscriberId
	p.Subscription.Subscriber = p.Subscriber
	p.Subscription.PlanID = p.ProcessData.PlanID
	p.Subscription.RemotePlanID = strconv.Itoa(p.PaymentReturn.RemotePlanID)
	p.Subscription.RemoteSubscriptionID = strconv.Itoa(p.PaymentReturn.RemoteSubscriptionID)
	p.Subscription.StartDate = startPeriod
	p.Subscription.ExpiresAt = endPeriod
	p.Subscription.Status = p.PaymentReturn.Status
	p.Subscription.CreatedAt = p.PaymentReturn.CreatedAt
	p.Subscription.UpdatedAt = p.PaymentReturn.UpdatedAt

	// Payment
	p.Payment.ID = paymentId
	p.Payment.TransactionID = strconv.Itoa(p.PaymentReturn.CurrentTransaction.RemoteTransactionID)
	p.Payment.Gateway = os.Getenv("GATEWAY")
	p.Payment.PaymentType = p.PaymentReturn.PaymentMethod
	p.Payment.Status = p.PaymentReturn.Status
	p.Payment.Total = p.PaymentReturn.CurrentTransaction.Amount
	p.Payment.CardBrand = p.PaymentReturn.CardBrand
	p.Payment.CardLastDigits = p.PaymentReturn.CardLastDigits
	p.Payment.BoletoURL = p.PaymentReturn.CurrentTransaction.BoletoURL
	p.Payment.BoletoBarcode = p.PaymentReturn.CurrentTransaction.BoletoBarcode
	p.Payment.BoletoExpirationDate = p.PaymentReturn.CurrentTransaction.BoletoExpirationDate
	p.Payment.Installments = p.PaymentReturn.CurrentTransaction.Installments
	p.Payment.SubscriptionID = subscriptionId
	p.Payment.CreatedAt = p.PaymentReturn.CreatedAt
	p.Payment.UpdatedAt = p.PaymentReturn.UpdatedAt

	// Subscriber
	p.Subscriber.ID = subscriberId
	p.Subscriber.Name = p.ProcessData.Name
	p.Subscriber.Email = p.ProcessData.Email
	p.Subscriber.DocumentNumber = p.ProcessData.DocumentNumber
	p.Subscriber.Street = p.ProcessData.Street
	p.Subscriber.StreetNumber = p.ProcessData.StreetNumber
	p.Subscriber.Complementary = p.ProcessData.Complementary
	p.Subscriber.Neighborhood = p.ProcessData.Neighborhood
	p.Subscriber.Zipcode = p.ProcessData.Zipcode
	p.Subscriber.DDD = p.ProcessData.DDD
	p.Subscriber.Number = p.ProcessData.PhoneNumber
	p.Subscriber.CreatedAt = p.PaymentReturn.CreatedAt
	p.Subscriber.UpdatedAt = p.PaymentReturn.UpdatedAt
	p.Subscriber.Subscriptions = models.Subscriptions{p.Subscription}

	p.Connection.ValidateAndCreate(&p.Subscriber)
	p.Connection.ValidateAndCreate(&p.Subscription)
	p.Connection.ValidateAndCreate(&p.Payment)

	return nil
}
