package dal

type AppDAL struct {
	User         UserDAL
	Token        TokenDAL
	Setting      SettingsDAL
	Media        MediaDAL
	Userdata     UserdataDAL
	Subscription SubscriptionDAL
	Postlist     PostlistDAL
	Tiers        TiersDAL
	Wallet       WalletDAL
	CryptoPrice  CryptoPriceDAL
	Transaction  TransactionDAL
}
