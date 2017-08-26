package harvest

import (
	"fmt"
	"net/http"
)

func (h *Harvest) LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := h.oauthConfig.AuthCodeURL("q")
	http.Redirect(w, r, url, 302)
}

func (h *Harvest) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")

	token, err := h.oauthConfig.Exchange(r.Context(), authCode)
	if err != nil {
		fmt.Println(err)
	}
	h.config.Token = token

	writeConfig(h.config)
}
