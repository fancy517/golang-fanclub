package types

type BaseFilter struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (f *BaseFilter) SetDefault() {
	if f.Page == 0 {
		f.Page = 1
	}

	if f.PageSize <= 0 {
		f.PageSize = 100
	}
}

type MatchFilter struct {
	BaseFilter
	Title    string `json:"title" form:"title"`
	Category string `json:"category" form:"category"`
	TeamName string `json:"team_name" form:"team_name"`
	LeagueID int    `json:"league_id" form:"league_id"`
	Status   string `json:"status" form:"status"` // live or all
	Date     string `json:"date" form:"date"`
	Sort     string `json:"sort" form:"sort"` // asc or desc by close_at field
}

type UserFilter struct {
	BaseFilter
	Query string `json:"query" form:"query"`
}

type UserMatchFilter struct {
	BaseFilter
	Type string `json:"type" form:"type"` // "live", "closed"
}

type TicketFilter struct {
	BaseFilter
}

type TxFilter struct {
	BaseFilter
	Type  string `json:"type" form:"type"`
	Year  int    `json:"year" form:"year"`
	Month int    `json:"month" form:"month"`
}

type SubscriberFilter struct {
	BaseFilter
}

type NewsletterFilter struct {
	BaseFilter
	Title string `json:"title" form:"title"`
}
