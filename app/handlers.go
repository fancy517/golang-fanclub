package app

import (
	"fanclub/internal/dal"
	adminhandler "fanclub/internal/handlers/admin"
	userhandler "fanclub/internal/handlers/user"
	wallethandler "fanclub/internal/handlers/wallet"
	cryptojob "fanclub/internal/jobs/crypto"
	"fanclub/internal/mailer"
)

const (
	MailQueueCapacity = 100
)

func (app *application) initHandlers() {
	// init mailer before initialize handlers
	app.mailer = mailer.New(
		app.config.smtp.host,
		app.config.smtp.port,
		app.config.smtp.username,
		app.config.smtp.password,
		app.config.smtp.sender,
		app.config.siteContactEmail,
		MailQueueCapacity,
	)

	// init dals
	appDal := dal.AppDAL{
		User:         dal.NewUserDAL(app.db),
		Token:        dal.NewTokenDAL(app.db),
		Setting:      dal.NewSettingsDAL(app.db),
		Media:        dal.NewMediaDAL(app.db),
		Userdata:     dal.NewUserdataDAL(app.db),
		Subscription: dal.NewSubscriptionDAL(app.db),
		Postlist:     dal.NewPostlistDAL(app.db),
		Tiers:        dal.NewTiersDAL(app.db),
		Wallet:       dal.NewWalletDAL(app.db),
		CryptoPrice:  dal.NewCryptoPriceDAL(app.db),
		Transaction:  dal.NewTransactionDAL(app.db),
	}

	// Init handlers
	app.userHandler = userhandler.NewHandler(appDal, app.mailer, app.troner)
	app.adminHandler = adminhandler.NewHandler(appDal, app.mailer)
	app.walletHandler = wallethandler.NewHandler(appDal, app.troner)

	// Init Jobs
	app.cryptoPriceJob = cryptojob.NewCryptoPriceCollector(appDal)
}
