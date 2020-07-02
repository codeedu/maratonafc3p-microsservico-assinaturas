package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"net/http"
	"os"
	"subscription_service/models"
	"subscription_service/services"
)


// SubscribeIndex default implementation.
func SubscribeIndex(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)

	// Allocate an empty Plan
	plan := &models.Plan{}

	// To find the Plan the parameter plan_id is used.
	if err := tx.Find(plan, c.Param("plan_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("GATEWAY_ENCRYPTION_KEY", os.Getenv("GATEWAY_ENCRYPTION_KEY"))
	c.Set("plan", plan)

	return c.Render(http.StatusOK, r.HTML("subscribe/index.html"))
}

// Process the subscription
func SubscribeProcess(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)

	processData := &services.ProcessData{}

	// Bind process to the html form elements
	if err := c.Bind(processData); err != nil {
		fmt.Println("====== error ====")
		return err
	}

	service := services.NewPaymentService()
	service.Connection = tx
	service.RabbitMQ = RabbitMQ
	err := service.Process(*processData)

	c.Set("GATEWAY_ENCRYPTION_KEY", os.Getenv("GATEWAY_ENCRYPTION_KEY"))
	if err != nil {
		c.Flash().Add("Declined","Transação negada. Tente novamente.")
		// Allocate an empty Plan
		plan := &models.Plan{}

		// To find the Plan the parameter plan_id is used.
		if err := tx.Find(plan, c.Param("plan_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}

		c.Set("plan", plan)
		return c.Render(http.StatusOK, r.HTML("subscribe/index.html"))
	}

	if service.PaymentReturn.Status == "unpaid" {
		c.Set("boletoURL", service.PaymentReturn.CurrentTransaction.BoletoURL)
		return c.Render(http.StatusOK, r.HTML("subscribe/boleto.html"))
	}

	return c.Render(http.StatusOK, r.HTML("subscribe/success.html"))
}
