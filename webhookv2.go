package telebot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//WebhookV2 enables webhook
type WebhookV2 struct {
	PublicURL      string
	AllowedUpdates []string
	bot            *Bot
}

//Register webhook
func (h *WebhookV2) Register(b *Bot) error {
	data := struct {
		URL            string   `json:"url"`
		AllowedUpdates []string `json:"allowed_updates"`
	}{
		h.PublicURL,
		h.AllowedUpdates,
	}

	res, err := b.Raw("setWebhook", data)
	if err != nil {
		b.debug(fmt.Errorf("setWebhook failed %q: %v", string(res), err))
		return err
	}
	var result registerResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		b.debug(fmt.Errorf("bad json data %q: %v", string(res), err))
		return err
	}
	if !result.Ok {
		err := fmt.Errorf("cannot register webhook: %s", result.Description)
		b.debug(err)
		return err
	}
	h.bot = b
	return nil
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
	h.bot.incomingUpdate(&update)
}
