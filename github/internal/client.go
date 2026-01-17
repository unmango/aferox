package internal

import (
	"net/http"
	"os"

	"github.com/google/go-github/v81/github"
)

func DefaultClient() *github.Client {
	client := github.NewClient(http.DefaultClient)

	if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client = client.WithAuthToken(token)
	}

	return client
}
