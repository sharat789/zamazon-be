package payment

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"log"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error)
	GetPaymentStatus(paymentId string) (*stripe.CheckoutSession, error)
}

type payment struct {
	apiKey     string
	successUrl string
	cancelURL  string
}

func (p payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.apiKey
	amountInCents := int64(amount * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				UnitAmount: stripe.Int64(amountInCents),
				Currency:   stripe.String(string(stripe.CurrencyEUR)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Electronics"),
				},
			},
			Quantity: stripe.Int64(1),
		},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(p.successUrl),
		CancelURL:  stripe.String(p.cancelURL),
	}

	params.AddMetadata("userId", fmt.Sprintf("%d", userId))
	params.AddMetadata("orderId", fmt.Sprintf("%s", orderId))

	session, err := session.New(params)
	if err != nil {
		log.Printf("Error while creating payment session %v\n", err.Error())
		return nil, errors.New("could not create payment session")
	}
	return session, nil
}

func (p payment) GetPaymentStatus(paymentId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.apiKey
	session, err := session.Get(paymentId, nil)
	if err != nil {
		log.Printf("Error while fetching payment status %v\n", err.Error())
		return nil, errors.New("could not fetch payment status")
	}
	return session, nil
}

func NewPaymentClient(apiKey, successUrl, failureUrl string) PaymentClient {
	return &payment{
		apiKey:     apiKey,
		successUrl: successUrl,
		cancelURL:  failureUrl,
	}
}
