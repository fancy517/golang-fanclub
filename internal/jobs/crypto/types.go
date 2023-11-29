package crypto

type usdvalue struct {
	USD float64 `json:"usd"`
}

type CoingeckoResponse struct {
	Tron usdvalue `json:"tron"`
}

type DiaDataResponse struct {
	Tron float64 `json:"price"`
}
