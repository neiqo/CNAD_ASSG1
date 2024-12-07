package models

type MemberBenefit struct {
	BenefitID      int    `json:"benefitID"`
	MembershipTier string `json:"membershipTier"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

type Promotion struct {
	PromotionID  int     `json:"promotionID"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Discount     float64 `json:"discount"`
	IfPercentage bool    `json:"ifPercentage"`
}
