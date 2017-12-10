package gui

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalletSpendHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/wallet/spend", nil)
	if err != nil {
		t.Fatal(err)

		rr := httptest.NewRecorder()
		http.HandlerFunc
		handler := http.HandlerFunc(WalletSpendHandler)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	gateway := &FakeGateway{}
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WalletSpendHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
