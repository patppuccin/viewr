package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/patppuccin/viewr/src/config"
	"github.com/patppuccin/viewr/src/constants"
	"github.com/patppuccin/viewr/src/helpers"
	"github.com/patppuccin/viewr/src/models"
	"github.com/patppuccin/viewr/src/out"
)

func Run(ctx context.Context, logLevel, address string, port int, logToConsole bool) error {
	// Set start time for logging
	startTime := time.Now()

	// Validate availability of address and port
	if err := helpers.CheckTCPBind(address, port); err != nil {
		return helpers.SafeErr("unable to bind to "+address+":"+strconv.Itoa(port), err)
	}

	// Set up signal handling for graceful shutdown
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
	}

	errChan := make(chan error, 1)

	// Prep for Server Context Step 1: Initialize logger
	logger, err := out.NewStructuredLogger(logLevel, logToConsole)
	if err != nil {
		return helpers.SafeErr("error initializing logger", err)
	}

	if logToConsole {
		os.Stdout.WriteString(out.Banner(constants.AppDescription) + "\n")
	}

	logger.Info().Msgf("initializing %s v%s", constants.AppFullName, constants.AppVersion)
	logger.Info().Msgf("configuration source: %s", config.GlobalConfigSrc)

	// Assemble server context
	serverCtx := &models.AppContext{
		Config: config.GlobalConfig,
		Logger: logger,
	}

	// Setup router
	router, err := setupRoutes(serverCtx)
	if err != nil {
		return helpers.SafeErr("error setting up routes", err)
	}

	// Setup server
	serveAddr := address + ":" + strconv.Itoa(port)
	server := &http.Server{
		Addr:    serveAddr,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		logger.Info().Msg("server is running at " + serveAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Warn().Msg("initiating server shutdown")
	case err := <-errChan:
		logger.Error().Msg("server encountered an error: " + err.Error())
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Msg("error during shutdown: " + err.Error())
	}
	// if err := dbConn.Close(); err != nil {
	// 	logger.Error().Msg("error closing database connection " + err.Error())
	// }

	logger.Info().Msgf("%s server shut down after %s", constants.AppAbbrName, time.Since(startTime).String())
	return nil
}
