// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"telegram-bot/internal/handler/operations"
	"telegram-bot/internal/worker"
	"telegram-bot/pkg/smap"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	MSGS_BUFFER  = 10
	TELBOT_TOKEN = "xxx"
)

//go:generate swagger generate server --target ../../../telegram-bot --name TelegramBot --spec ../../api/api.yml --server-package internal/handler

func configureFlags(api *operations.TelegramBotAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TelegramBotAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	msgs := make(chan string, MSGS_BUFFER)
	users := smap.NewSTMap(100)
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		panic(err)
	}
	go worker.InitTelBotConsumer(users, bot)
	go worker.InitTelBotProducer(users, bot, msgs)

	api.PostSendHandler = operations.PostSendHandlerFunc(func(params operations.PostSendParams) middleware.Responder {
		err := SendMsg(msgs, params.Msg)
		if err != nil {
			return operations.NewPostSendInternalServerError()
		}
		return operations.NewPostSendOK()
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
