package bankid

import (
	"context"
	"github.com/NicklasWallgren/bankid/configuration"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

// BankId contains the validator and configuration context
type BankId struct {
	validator     *validator.Validate
	configuration *configuration.Configuration
	client        *Client
}

// NewBankId returns a new instance of 'BankId'
func NewBankId(configuration *configuration.Configuration) *BankId {
	return &BankId{validator: newValidator(), configuration: configuration}
}

// Authenticate - Initiates an authentication order.
//
// Use the collect method to query the status of the order.
// If the request is successful, the orderRef and autoStartToken is returned.
func (b BankId) Authenticate(payload *AuthenticationPayload) (*AuthenticateResponse, error) {
	request := newAuthenticationRequest(payload)
	response, err := b.call(request)

	if err != nil {
		return nil, err
	}

	authenticateResponse := (*response).(*AuthenticateResponse)
	return authenticateResponse, nil
}

// Sign - Initiates an sign order.
//
// Use the collect method to query the status of the order.
// If the request is successful, the orderRef and autoStartToken is returned.
func (b BankId) Sign(payload *SignPayload) (*SignResponse, error) {
	request := newSignRequest(payload)
	response, err := b.call(request)

	if err != nil {
		return nil, err
	}

	signResponse := (*response).(*SignResponse)
	return signResponse, nil
}

// Collect - Collects the result of a sign or auth order suing the orderRef as reference.
//
// RP should keep calling collect every two seconds as long as status indicates pending.
// RP must abort if status indicates failed. The user identity is returned when complete.
func (b BankId) Collect(payload *CollectPayload) (*CollectResponse, error) {
	request := newCollectRequest(payload)
	response, err := b.call(request)

	if err != nil {
		return nil, err
	}

	collectResponse := (*response).(*CollectResponse)
	return collectResponse, nil

}

// Cancel - Cancels an ongoing sign or auth order.
//
// This is typically used if the user cancels the order in your service or app.
func (b BankId) Cancel(payload *CancelPayload) (*CancelResponse, error) {
	request := newCancelRequest(payload)
	response, err := b.call(request)

	if err != nil {
		return nil, err
	}

	cancelResponse := (*response).(*CancelResponse)
	return cancelResponse, nil
}

// call validates the prerequisites of the requests and invokes the REST API method
func (b *BankId) call(request Request) (*Response, error) {
	context, cancel := context.WithTimeout(context.Background(), b.configuration.Timeout*time.Second)
	defer cancel()

	// Validate the integrity of the call
	err := b.validator.Struct(request.Payload())

	if err != nil {
		return nil, err
	}

	if err = b.initialize(); err != nil {
		return nil, err
	}

	response, err := (*b.client).call(request, &context, b)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// initialize prepares the client in head of a request
func (b *BankId) initialize() error {
	// Check whether the client has been initialized
	if b.client != nil {
		return nil
	}

	// Lazy initialization
	client, err := newClient(b.configuration)

	if err != nil {
		return err
	}

	b.client = &client

	return nil
}
