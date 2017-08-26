package harvest

import (
	"os"

	"github.com/bram96/timeular-go/zei"

	"golang.org/x/oauth2"
)

type Side struct {
	ProjectID string `yaml:"project_id"`
	TaskID    string `yaml:"task_id"`
}

type Config struct {
	Sides map[zei.Position]Side `yaml:"sides"`
	Token *oauth2.Token         `yaml:"token"`
}

type Harvest struct {
	config      Config
	oauthConfig oauth2.Config

	currentRunningTimer int
}

func NewHarvest() *Harvest {
	c, err := readConfig()
	if err != nil {
		panic(err)
	}
	return &Harvest{
		config: c,
		oauthConfig: oauth2.Config{
			ClientID:     os.Getenv("HARVEST_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("HARVEST_OAUTH_CLIENT_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.harvestapp.com/oauth2/authorize",
				TokenURL: "https://api.harvestapp.com/oauth2/token",
			},
			RedirectURL: "http://localhost:8080/harvest/callback",
		},
	}
}
