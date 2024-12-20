package main

import (
	"flag"
	"fmt"

	"github.com/dexhealth/fhir/pkg/core/meta"
	"github.com/dexhealth/fhir/pkg/core/search/db_funcs"
	"github.com/dexhealth/fhir/pkg/core/search/lookup"
	"github.com/dexhealth/fhir/pkg/r4/datastore"
	"github.com/dexhealth/fhir/pkg/r4/datastore/mongodb"
	models "github.com/dexhealth/fhir/pkg/r4/resource"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

const (
	SEARCH_PARAM_FILE = "search-parameters.json"
)

func main() {
	println("hello")

	lookupTable, _ := lookup.NewSearchTable(SEARCH_PARAM_FILE)
	metaCtrl := meta.NewController(lookupTable)

	// args

	portFlag := flag.String("port", "3000", "http server port")
	urlFlag := flag.String("url", "http://localhost:3000", "public api url")
	flag.Parse()

	// logger

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	// router

	e := echo.New()

	// bindings

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogLatency:       true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.String("method", v.Method),
				zap.String("request_id", v.RequestID),
				zap.String("remote_ip", v.RemoteIP),
				zap.Duration("latency", v.Latency),
				zap.String("bytes_in", v.ContentLength),
				zap.Int64("bytes_out", v.ResponseSize),
			)
			return nil
		},
	}))

	var ds datastore.IDatastore
	var err error
	db_funcs.Generate(SEARCH_PARAM_FILE)
	ds, err = mongodb.NewDatastore()
	if err != nil {
		fmt.Println("db creation error")
		return
	}

	// routes
	e.POST("/", metaCtrl.Post)
	models.UseRoutes(e, ds, *urlFlag, lookupTable)

	e.Logger.Fatal(e.Start(":" + *portFlag))
}
