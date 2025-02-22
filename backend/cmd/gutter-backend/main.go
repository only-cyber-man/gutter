package main

import (
	"errors"
	"log/slog"

	"gutter/internal/handler"

	"github.com/tomek7667/cyberman-go/crypto"
	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/logger"
	"github.com/tomek7667/cyberman-go/pocketbase"
	"github.com/tomek7667/cyberman-go/rest"
	"github.com/tomek7667/cyberman-go/utils"
)

var (
	pocketbaseClient *pocketbase.Client
	restClient       *rest.Client
	expoClient       *expo.Client
)

func init() {
	var errs, err error
	slog.Info("Initializing logger")
	logger.SetLogLevel()

	APP_PORT := utils.Getenv("APP_PORT", "7005")
	restClient = rest.New(APP_PORT)

	PB_URL := utils.Getenv("PB_URL", "http://127.0.0.1:8090")
	PB_USERNAME := utils.Getenv("PB_USERNAME", "local@admin.admin")
	PB_PASSWORD := utils.Getenv("PB_PASSWORD", "adminadmin")
	pocketbaseClient, err = pocketbase.New(pocketbase.ClientInput{
		Url:      PB_URL,
		Username: PB_USERNAME,
		Password: PB_PASSWORD,
	})
	if err != nil {
		errs = errors.Join(errs, err)
	}
	err = pocketbaseClient.RefreshToken()
	if err != nil {
		errs = errors.Join(errs, err)
	}

	expoToken := utils.Getenv("EXPO_ACCESS_TOKEN", "paste-expo-access-token-here")
	slog.Info(
		"expo access token loaded",
		"EXPO_ACCESS_TOKEN", crypto.Obfuscate(expoToken, 8),
	)
	expoClient = expo.New(expoToken)

	if errs != nil {
		slog.Error("Failed to initialize gutter", "error", errs.Error())
		panic("Failed to initialize the gutter")
	}
	slog.Info(
		"jwt loaded",
		"JWT_SECRET", crypto.Obfuscate(utils.Getenv(
			"JWT_SECRET",
			"8c7fafb856380624fa60b22e7baf311d",
		), 8),
	)
	slog.Info("Connected")
}

func main() {
	h := handler.New(
		restClient,
		pocketbaseClient,
		expoClient,
	)
	h.Start()
}
