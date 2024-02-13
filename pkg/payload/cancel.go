package payload

import "github.com/NicklasWallgren/bankid/v2/pkg/internal/http"

// CancelPayload holds the required fields of the collect Payload.
type CancelPayload struct {
	http.Payload `json:"-"`
	// The orderRef from the response from authentication or sign.
	OrderRef string `json:"orderRef"`
}