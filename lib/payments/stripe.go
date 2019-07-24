package payments

import (
	"errors"
	"fmt"

	"github.com/jasongwartz/carbon-offset-backend/lib/config"
	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"

	stripe "github.com/stripe/stripe-go"
	paymentintent "github.com/stripe/stripe-go/paymentintent"
)

func init() {
	// Initialize stripe-go with Stripe secret key from environment
	stripe.Key = config.C.StripeSecretKey
}

func generatePaymentResponse(intent *stripe.PaymentIntent) (*models.PaymentResponse, error) {
	if intent.Status == stripe.PaymentIntentStatusRequiresAction &&
		intent.NextAction.Type == "use_stripe_sdk" {

		return &models.PaymentResponse{
			RequiresAction:            stripe.Bool(true),
			PaymentIntentClientSecret: &intent.ClientSecret,
		}, nil
	} else if intent.Status == stripe.PaymentIntentStatusSucceeded {
		return &models.PaymentResponse{
			Success: stripe.Bool(true),
		}, nil
	}
	return &models.PaymentResponse{}, fmt.Errorf("Invalid Payment Intent status: %s", intent.Status)

}

func Checkout(paymentMethod string, amount int, currency models.Currency) (*models.PaymentResponse, error) {
	params := &stripe.PaymentIntentParams{
		PaymentMethod: stripe.String(paymentMethod),
		Amount:        stripe.Int64(int64(amount)),
		Currency:      stripe.String(string(currency)),
		ConfirmationMethod: stripe.String(string(
			stripe.PaymentIntentConfirmationMethodManual,
		)),
		Confirm: stripe.Bool(true),
	}
	intent, err := paymentintent.New(params)

	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			// Display error on client
			return nil, errors.New(stripeErr.Msg)
		}
		return nil, err // TODO: log server-side, show generic message to client
	}

	return generatePaymentResponse(intent)
}

func Confirm(paymentIntent string) (*models.PaymentResponse, error) {
	intent, err := paymentintent.Confirm(
		paymentIntent, nil,
	)

	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			// Display error on client
			return nil, errors.New(stripeErr.Msg)
		}
		return nil, err // TODO: handle this error server-side, return client-friendly error
	}

	return generatePaymentResponse(intent)
}
