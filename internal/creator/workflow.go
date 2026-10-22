package creator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Subscriber struct {
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
}
type Asset struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	DownloadURL string `json:"download_url"`
}
type Alert struct{ Phone, Message, RequestID string }

func ProcessDelivery(sub Subscriber, asset Asset) (Alert, error) {
	if !sub.Active {
		return Alert{}, fmt.Errorf("subscriber is inactive")
	}
	if sub.Phone == "" || asset.ID == "" || asset.DownloadURL == "" {
		return Alert{}, fmt.Errorf("subscriber and asset details are required")
	}
	title := strings.TrimSpace(asset.Title)
	if title == "" {
		title = "your digital item"
	}
	sum := sha256.Sum256([]byte(asset.ID + "|" + sub.Phone))
	id := hex.EncodeToString(sum[:8])
	return Alert{Phone: sub.Phone, Message: fmt.Sprintf("Delivery ready: %s %s", title, asset.DownloadURL), RequestID: "delivery-" + id}, nil
}
