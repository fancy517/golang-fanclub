package app

import (
	"database/sql"
	adminhandler "fanclub/internal/handlers/admin"
	userhandler "fanclub/internal/handlers/user"
	wallethandler "fanclub/internal/handlers/wallet"
	"fanclub/internal/mailer"
	"log"

	"fanclub/internal/interfaces"
	"fanclub/internal/troner"

	"github.com/gin-gonic/gin"
)

type application struct {
	engine *gin.Engine
	config AppConfig
	db     *sql.DB

	// Handlers
	userHandler   userhandler.Handler
	adminHandler  adminhandler.Handler
	walletHandler wallethandler.Handler
	// Block Scanner
	// scanner scanner.BlockScanner

	// Mailer
	mailer         mailer.Mailer
	cryptoPriceJob interfaces.CryptoPriceCollector
	troner         interfaces.Troner
}

func NewApp() (interfaces.App, error) {
	app := application{}

	// first of all, load config from .env
	if err := app.loadConfig(); err != nil {
		return nil, err
	}

	// Open database
	if err := app.openDB(); err != nil {
		return nil, err
	}

	host := troner.Mainnet
	if app.config.chain == "Nile" {
		host = troner.NileTestnet
	} else if app.config.chain == "Shasta" {
		host = troner.ShastaTestnet
	}
	app.troner = troner.NewTroner(
		host,
		app.config.house.address,
		app.config.house.privateKey,
		app.config.house.minDeposit,
		app.config.house.earningFee,
	)

	app.initHandlers()
	app.initRouters()

	return &app, nil
}

func (app *application) Serve() {
	// app.scanner.Run()
	app.mailer.Run()
	app.cryptoPriceJob.Run()

	if err := app.engine.Run(); err != nil {
		panic(err)
	}
}

func (app *application) CloseDB() {
	app.db.Close()
	log.Println("Closing database connection")
}
