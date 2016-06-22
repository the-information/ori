package models

import (
	"strings"
	"time"
)

type Offer struct {
	CreatedAt          time.Time `json:"createdAt"`
	LastUpdatedAt      time.Time `json:"lastUpdatedAt"`
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ValidFrom          time.Time `json:"validFrom"`
	ValidUntil         time.Time `json:"validUntil"`
	ValidClaimants     []string  `json:"validClaimants"`
	ValidEmailSuffix   string    `json:"validEmailSuffix"`
	ClaimCount         int       `json:"claimCount"`
	ClaimsRemaining    int       `json:"claimsRemaining"`
	BillingCycleLength string    `json:"billingCycleLength"`
	Price              string    `json:"price"`
	PromotionalCycles  int       `json:"promotionalCycles"`
	PromotionalPrice   string    `json:"promotionalPrice"`
}

func (o *Offer) ClaimableBy(email string) bool {

	for _, validClaimant := range o.ValidClaimants {
		if email == validClaimant {
			return true
		}
	}

	return strings.HasSuffix(email, o.ValidEmailSuffix)

}
