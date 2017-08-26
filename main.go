package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bram96/timeular-go/harvest"
	"github.com/bram96/timeular-go/zei"
	"github.com/gorilla/mux"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	h := harvest.NewHarvest()
	z, err := zei.NewZEIConnection(h)
	if err != nil {
		log.Fatal(err)
	}

	defer z.Disconnect()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json, _ := json.Marshal(struct {
			Status   zei.Status `json:"status"`
			Position string     `json:"position"`
		}{z.ConnectionStatus(), zei.PositionName[z.Position()]})
		w.Write(json)
	})

	harvestRouter := mux.NewRouter()
	harvestRouter.HandleFunc("/harvest/login", h.LoginHandler)
	harvestRouter.HandleFunc("/harvest/callback", h.CallbackHandler)
	http.Handle("/harvest/", harvestRouter)

	http.ListenAndServe(":8080", nil)

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(z.Position())
	}
}
