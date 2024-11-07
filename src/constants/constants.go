package constants

import (
	"github.com/RubenPari/clear-songs/environment"
	"os"
)

var clientId = os.Getenv("CLIENT_ID")

func GetClientId() string {
	if clientId == "" {
		environment.LoadEnvVariables()
	}

	return clientId
}

var clientSecret = os.Getenv("CLIENT_SECRET")

func GetClientSecret() string {
	if clientSecret == "" {
		environment.LoadEnvVariables()
	}

	return clientSecret
}

var redirectUrl = os.Getenv("REDIRECT_URL")

func GetRedirectUrl() string {
	if redirectUrl == "" {
		environment.LoadEnvVariables()
	}

	return redirectUrl
}

var dbHost = os.Getenv("DB_HOST")

func GetDbHost() string {
	if dbHost == "" {
		environment.LoadEnvVariables()
	}

	return dbHost
}

var dbUser = os.Getenv("DB_USER")

func GetDbUser() string {
	if dbUser == "" {
		environment.LoadEnvVariables()
	}

	return dbUser
}

var dbPassword = os.Getenv("DB_PASSWORD")

func GetDbPassword() string {
	if dbPassword == "" {
		environment.LoadEnvVariables()
	}

	return dbPassword
}

var dbName = os.Getenv("DB_NAME")

func GetDbName() string {
	if dbName == "" {
		environment.LoadEnvVariables()
	}

	return dbName
}

var dbPort = os.Getenv("DB_PORT")

func GetDbPort() string {
	if dbPort == "" {
		environment.LoadEnvVariables()
	}

	return dbPort
}

var Scopes = []string{
	"user-read-private",
	"user-read-email",
	"user-library-read",
	"user-library-modify",
	"playlist-read-private",
	"playlist-read-collaborative",
	"playlist-modify-public",
	"playlist-modify-private",
}

const LimitGetPlaylistTracks = 100
const LimitRemovePlaylistTracks = 100
const Offset = 0
