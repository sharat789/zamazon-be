package payment

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"log"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(paymentId string) (*stripe.PaymentIntent, error)
	CreateCheckoutSession(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error)
	GetCheckoutSession(sessionId string) (*stripe.CheckoutSession, error)
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

func (p payment) CreateCheckoutSession(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.apiKey
	amountInCents := int64(amount * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyEUR)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Order Payment"),
					},
					UnitAmount: stripe.Int64(amountInCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(p.successUrl + "/#/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(p.cancelURL + "/#/cart"),
	}

	params.AddMetadata("userId", fmt.Sprintf("%d", userId))
	params.AddMetadata("orderId", fmt.Sprintf("%s", orderId))

	session, err := session.New(params)
	if err != nil {
		log.Printf("Error while creating checkout session %v\n", err.Error())
		return nil, errors.New("could not create checkout session")
	}

	return session, nil
}

func (p payment) GetCheckoutSession(sessionId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.apiKey
	session, err := session.Get(sessionId, nil)
	if err != nil {
		log.Printf("Error retrieving checkout session: %v\n", err)
		return nil, errors.New("could not retrieve checkout session")
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
