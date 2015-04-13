package app

import (
	"log"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/Altoros/cf-cassandra-service-broker/api"
	"github.com/Altoros/cf-cassandra-service-broker/config"
)

type AppContext struct {
	appConfig  *config.Config
	negroni    *negroni.Negroni
	mainRouter *mux.Router
}

func NewApp(appConfig *config.Config) (*AppContext, error) {
	var err error

	app := new(AppContext)
	app.appConfig = appConfig

	negroni := negroni.Classic()
	mainRouter := mux.NewRouter()
	negroni.UseHandler(mainRouter)
	app.negroni = negroni
	app.mainRouter = mainRouter

	err = app.initApi()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *AppContext) Start() {
	log.Println("Starting broker on port", app.appConfig.PortStr())
	app.negroni.Run(":" + app.appConfig.PortStr())
}

func (app *AppContext) Stop() {
	log.Println("Stopping broker")
	os.Exit(0)
}

func (app *AppContext) initApi() error {
	apiRouter := app.mainRouter.PathPrefix("/v2").Subrouter()

	catalogController := api.NewCatalogController(&app.appConfig.Catalog)
	catalogController.AddRoutes(apiRouter)

	return nil
}