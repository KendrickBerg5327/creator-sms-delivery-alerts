package creator

import (
	"strings"
	"testing"
)

func TestProcessDeliveryDecision(t *testing.T) {
	tests := []struct {
		name    string
		sub     Subscriber
		asset   Asset
		wantErr bool
	}{
		{"active subscriber receives link", Subscriber{Phone: "+15551234567", Active: true}, Asset{ID: "pack-7", Title: "March pack", DownloadURL: "https://cdn.example/pack-7"}, false},
		{"inactive subscriber is skipped", Subscriber{Phone: "+15551234567", Active: false}, Asset{ID: "pack-7", DownloadURL: "https://cdn.example/pack-7"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, e := ProcessDelivery(tt.sub, tt.asset)
			if (e != nil) != tt.wantErr {
				t.Fatalf("error=%v", e)
			}
			if e == nil && (!strings.Contains(a.Message, "https://cdn.example/pack-7") || a.RequestID == "") {
				t.Fatalf("bad alert: %+v", a)
			}
		})
	}
}
