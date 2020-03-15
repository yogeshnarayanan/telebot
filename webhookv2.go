package telebot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//WebhookV2 enables webhook
type WebhookV2 struct {
	dest chan<- Update
	bot  *Bot
}

//Poll poller
func (h *WebhookV2) Poll(b *Bot, dest chan Update, stop chan struct{}) {
	h.dest = dest
	h.bot = b

	go func(stop chan struct{}) {
		<-stop
		close(stop)
	}(stop)
}

// ServeHTTP The handler simply reads the update from the body of the requests
// and writes them to the update channel.
func (h *WebhookV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		h.bot.debug(fmt.Errorf("cannot decode update: %v", err))
		return
	}
	h.dest <- update
}
