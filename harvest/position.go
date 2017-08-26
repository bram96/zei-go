package harvest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bram96/timeular-go/zei"
)

func (h *Harvest) OnPositionChanged(p zei.Position) {
	fmt.Printf("Position changed: %v\n", p)
	c := h.oauthConfig.Client(context.Background(), h.config.Token)

	dat, ok := h.config.Sides[p]
	if ok {

		reqBody, _ := json.Marshal(struct {
			ProjectID string `json:"project_id"`
			TaskID    string `json:"task_id"`
			SpentAt   string `json:"spent_at"`
		}{dat.ProjectID, dat.TaskID, time.Now().Format("2006-01-02")})

		req, err := http.NewRequest("POST", "https://api.harvestapp.com/daily/add", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		r, err := c.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var currentTimer struct {
			ID int `json:"id"`
		}

		if err := json.Unmarshal(responseData, &currentTimer); err != nil {
			fmt.Println(err)
		}
		h.currentRunningTimer = currentTimer.ID
	} else if h.currentRunningTimer != 0 {
		//TODO stop timer

		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.harvestapp.com/daily/timer/%v", h.currentRunningTimer), bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		r, err := c.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(responseData))
	}
}
