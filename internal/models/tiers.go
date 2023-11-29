package models

type TierPlan struct {
	Duration string `json:"duration"`
	Price    string `json:"price"`
}

type TierModel struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Title     string `json:"tier_name"`
	Color     string `json:"tier_color"`
	Benefit   string `json:"tier_benefit"`
	Child     string `json:"tier_child"`
	Baseprice string `json:"base_price"`
	Month2    string `json:"month_two"`
	Month3    string `json:"month_three"`
	Month6    string `json:"month_six"`
	Active    string `json:"active"`
}

type Benefit struct {
	ID   string `json:"id"`
	Text string `json:"description"`
}

type RelativeTier struct {
	ID    string `json:"id"`
	Title string `json:"tier_name"`
}

type EditTierModel struct {
	Title         string  `json:"tier_name"`
	Color         string  `json:"tier_color"`
	Baseprice     string  `json:"base_price"`
	Benefits      *string `json:"tier_benefit"`
	RelativeTiers []RelativeTier
	Children      *string `json:"tier_child"`
	Month_two     string  `json:"month_two"`
	Month_three   string  `json:"month_three"`
	Month_six     string  `json:"month_six"`
	Active        string  `json:"active"`
}
