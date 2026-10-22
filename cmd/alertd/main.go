package main

import (
	"creator-sms-alerts/internal/creator"
	"creator-sms-alerts/internal/infrai"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type request struct {
	Subscriber creator.Subscriber `json:"subscriber"`
	Asset      creator.Asset      `json:"asset"`
}

func main() {
	c, e := infrai.New()
	if e != nil {
		log.Fatal(e)
	}
	http.HandleFunc("/delivery", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}
		var in request
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid json", 400)
			return
		}
		a, e := creator.ProcessDelivery(in.Subscriber, in.Asset)
		if e != nil {
			http.Error(w, e.Error(), 422)
			return
		}
		id, e := c.SendSMS(a.Phone, a.Message, a.RequestID)
		if e != nil {
			http.Error(w, e.Error(), 502)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message_id": id, "request_id": a.RequestID})
	})
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}
