//go:build integration

package api

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/mnty4/booking/model"
)

func NewUser() model.User {
	return model.User{
		Email:     "John.Doe@email.com",
		FirstName: "John",
		LastName:  "Doe",
	}
}

func TestUserPost(t *testing.T) {
	testCases := []struct {
		msg         string
		payloadJSON string
		code        int
	}{
		{"valid user", `{
			"email": "John.Doe@email.com",
			"firstName": "John",
			"lastName": "Doe"
		}`, 201},
		{"valid user", `{
			"email": "John.Doe@email.com",
			"lastName": "Doe"
		}`, 400},
		{"valid user", `{
			"email": "John.Doe@email.com",
			"firstName": "John"
		}`, 400},
		{"valid user", `{
			"firstName": "John",
			"lastName": "Doe"
		}`, 400},
	}
	for _, tc := range testCases {
		req, err := http.NewRequest("POST", "http://localhost:8080/api/users", bytes.NewReader([]byte(tc.payloadJSON)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != tc.code {
			t.Logf("expected status code %d but got %d.", tc.code, resp.StatusCode)
			t.Fail()
		}
	}
}
