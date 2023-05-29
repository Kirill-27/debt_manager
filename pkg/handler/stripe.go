package handler

import (
	"github.com/kirill-27/debt_manager/data"
	"log"
	"time"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
)

const stripeKey = "sk_test_51NCLXnL4aVSmf4QOQy6EZqDUt4f7ej8My8LFe462HwiWePr1MaAig6MHnGkmEksD10HOVJMUiNSfPZhO1UXSSsIY00HsiXDJXM"

func (h *Handler) StripeKeeper() {
	stripe.Key = stripeKey
	for {
		time.Sleep(20 * time.Second)
		lastTime, err := h.services.StripePaymentKeys.GetLastHandled()
		if err != nil {
			log.Println(err)
			continue
		}
		lastTimeMinusFive := lastTime.Add(-5 * time.Minute)

		params := &stripe.PaymentIntentListParams{
			CreatedRange: &stripe.RangeQueryParams{
				GreaterThan: lastTimeMinusFive.Unix(),
			},
		}
		var newTime *int64

		it := paymentintent.List(params)
		for it.Next() {

			paymentIntent := it.PaymentIntent()
			if newTime == nil {
				newTime = &paymentIntent.Created
			}
			stripePayment, err := h.services.StripePayment.GetStripePaymentByPaymentId(paymentIntent.ID)
			if stripePayment != nil {
				continue
			}

			newStripePayment := data.StripePayment{
				PaymentId: paymentIntent.ID,
			}

			if paymentIntent.Status != stripe.PaymentIntentStatusProcessing && paymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
				newStripePayment.UserId = 0
				newStripePayment.Status = data.StripePaymentsStatusCanceled
				_, err = h.services.StripePayment.CreateStripePayment(newStripePayment)
				if err != nil {
					log.Println(err)
					continue
				}
				continue
			}

			pm, err := paymentmethod.Get(paymentIntent.PaymentMethod.ID, nil)
			if err != nil {
				log.Println(err)
				continue
			}

			user, err := h.services.Authorization.GetUser(&pm.BillingDetails.Email, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			if user == nil {
				newStripePayment.Status = data.StripePaymentsStatusIncorrectEmail
				newStripePayment.UserId = 0
				_, err = h.services.StripePayment.CreateStripePayment(newStripePayment)
				if err != nil {
					log.Println(err)
					continue
				}
				continue
			}

			newStripePayment.UserId = user.Id
			newStripePayment.Status = convertStatus(paymentIntent.Status)
			if paymentIntent.Status == stripe.PaymentIntentStatusSucceeded {
				user.SubscriptionType = data.SubscriptionTypePremium
				err = h.services.Authorization.UpdateUser(*user)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			_, err = h.services.StripePayment.CreateStripePayment(newStripePayment)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		if err := it.Err(); err != nil {
			log.Println(err)
			continue

		}

		if newTime != nil {
			err = h.services.StripePaymentKeys.SetLastHandled(*newTime)
			if err != nil {
				log.Println(err)
			}
		}

	}
}

func (h *Handler) StripeHandler() {
	stripe.Key = stripeKey
	for {
		time.Sleep(20 * time.Second)
		status := data.StripePaymentsStatusProcessing
		processingPayments, err := h.services.StripePayment.GetAllStripePayments(&status, nil)
		if err != nil {
			log.Println(err)
		}
		for _, payment := range processingPayments {
			pi, err := paymentintent.Get(payment.PaymentId, nil)
			if err != nil {
				log.Println(err)
			}
			if pi.Status == stripe.PaymentIntentStatusProcessing {
				continue
			}
			if pi.Status == stripe.PaymentIntentStatusSucceeded {
				pm, err := paymentmethod.Get(pi.PaymentMethod.ID, nil)
				if err != nil {
					log.Println(err)
					continue
				}

				user, err := h.services.Authorization.GetUser(&pm.BillingDetails.Email, nil)
				user.SubscriptionType = data.SubscriptionTypePremium
				err = h.services.Authorization.UpdateUser(*user)
				if err != nil {
					log.Println(err)
					continue
				}

				err = h.services.StripePayment.UpdateStripePaymentStatus(payment.PaymentId, data.StripePaymentsStatusSucceeded)
				if err != nil {
					log.Println(err)
				}
				continue
			}
			err = h.services.StripePayment.UpdateStripePaymentStatus(payment.PaymentId, data.StripePaymentsStatusCanceled)
			if err != nil {
				log.Println(err)
			}

		}
	}
}

func convertStatus(paymentIntentStatus stripe.PaymentIntentStatus) int {
	if paymentIntentStatus == stripe.PaymentIntentStatusProcessing {
		return data.StripePaymentsStatusProcessing
	}
	if paymentIntentStatus == stripe.PaymentIntentStatusSucceeded {
		return data.StripePaymentsStatusSucceeded
	}
	return data.StripePaymentsStatusCanceled

}
