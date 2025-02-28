package payment

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"log"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(paymentId string) (*stripe.PaymentIntent, error)
}

type payment struct {
	apiKey     string
	successUrl string
	cancelURL  string
}

func (p payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.apiKey
	amountInCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
	}

	params.AddMetadata("userId", fmt.Sprintf("%d", userId))
	params.AddMetadata("orderId", fmt.Sprintf("%s", orderId))

	pi, err := paymentintent.New(params)

	if err != nil {
		log.Printf("Error while creating payment intent %v\n", err.Error())
		return nil, errors.New("could not create payment intent")
	}
	return pi, nil
}

func (p payment) GetPaymentStatus(paymentId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.apiKey
	params := &stripe.PaymentIntentParams{}
	result, err := paymentintent.Get(paymentId, params)
	if err != nil {
		log.Printf("Error while fetching payment status %v\n", err.Error())
		return nil, errors.New("could not fetch payment status")
	}
	return result, nil
}

func NewPaymentClient(apiKey, successUrl, failureUrl string) PaymentClient {
	return &payment{
		apiKey:     apiKey,
		successUrl: successUrl,
		cancelURL:  failureUrl,
	}
}
