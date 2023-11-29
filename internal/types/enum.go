package types

type BetDirection int
type MatchStatusType int
type UserStatusType int
type TxType string
type NotificationSourceType string

const (
	BetDirectionUnknown BetDirection = iota
	BetDirectionHome
	BetDirectionAway
	BetDirectionDraw
)

const (
	MatchStatusUnknown MatchStatusType = iota
	MatchStatusInit
	MatchStatusLive
	MatchStatusHomeWin
	MatchStatusAwayWin
	MatchStatusDraw
	MatchStatusCanceled
)

const (
	UserStatusUnverified UserStatusType = 0
	UserStatusVerified   UserStatusType = 1
	UserStatusBlocked    UserStatusType = -1
)

const (
	TxTypeDeposit    TxType = "D"
	TxTypeWithdrawal TxType = "W"
	TxTypeCollect    TxType = "C"
)

const (
	NotificationSourceMatchDecided  NotificationSourceType = "match_decided"
	NotificationSourceMatchRefunded NotificationSourceType = "match_refunded"
	NotificationSourceReferral      NotificationSourceType = "referral"
)
