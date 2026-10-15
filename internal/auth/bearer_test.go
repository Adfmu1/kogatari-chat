package auth

import (
	"net/http"
	"testing"
)

func TestBearer(t *testing.T) {
	headers := make([]http.Header, 3)

	headers[0] = http.Header{}
	headers[1] = http.Header{}
	headers[1].Add("Authorization", "Bearer")
	headers[2] = http.Header{}
	headers[2].Add("Authorization", "Bearer ok")

	for i, header := range headers {
		if i != 2 {
			_, err := GetBearerToken(header)
			if err == nil {
				t.Errorf("Expected error, got %v", err)
			}
		} else {
			_, err := GetBearerToken(header)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		}
	}
}
