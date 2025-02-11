//go:build integration

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mnty4/booking/errutil"
	"github.com/mnty4/booking/model"
)

func NewUser() model.User {
	return model.User{
		Email:     "John.Doe@email.com",
		FirstName: "John",
		LastName:  "Doe",
	}
}

func TestCreateUserValid(t *testing.T) {
	testCases := []struct {
		msg         string
		payloadJSON string
	}{
		{"happy path", `{
			"email": "John.Doe@email.com",
			"firstName": "John",
			"lastName": "Doe"
		}`},
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
		if resp.StatusCode != http.StatusCreated {
			t.Logf("expected status code %d but got %d.", http.StatusCreated, resp.StatusCode)
			t.Fail()
		}
	}
}

func TestCreateUserInvalid(t *testing.T) {
	testCases := []struct {
		msg         string
		payloadJSON string
		code        int
		status      string
	}{
		{"missing firstName", `{
			"email": "John.Doe@email.com",
			"lastName": "Doe"
		}`, 400, "VALIDATION"},
		{"missing lastName", `{
			"email": "John.Doe@email.com",
			"firstName": "John"
		}`, 400, "VALIDATION"},
		{"missing email", `{
			"firstName": "John",
			"lastName": "Doe"
		}`, 400, "VALIDATION"},
		{"empty firstName", `{
			"email": "John.Doe@email.com",
			"firstName": "",
			"lastName": "Doe"
		}`, 400, "VALIDATION"},
		{"empty lastName", `{
			"email": "John.Doe@email.com",
			"firstName": "John,
			"lastName": ""
		}`, 400, "VALIDATION"},
		{"email no @", `{
			"email": "johndoe",
			"firstName": "",
			"lastName": "Doe"
		}`, 400, "VALIDATION"},
		{"email no domain", `{
			"email": "johndoe@",
			"firstName": "",
			"lastName": "Doe"
		}`, 400, "VALIDATION"},
		{"email no local part", `{
			"email": "@johndoe",
			"firstName": "",
			"lastName": "Doe"
		}`, 400, "VALIDATION"},
	}
	client := new(http.Client)
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("POST", "http://localhost:8080/api/users", bytes.NewReader([]byte(tc.payloadJSON)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.code {
				t.Errorf("expected status code %d but got %d.", tc.code, resp.StatusCode)
			}
			var apiErr errutil.APIError
			if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
				t.Fatal(err)
			}
			if apiErr.Code != tc.code {
				t.Errorf("expected status code %d but got %d.", tc.code, apiErr.Code)
			}
			if apiErr.GetStatus() != tc.status {
				t.Errorf("expected status %q but got %q.", tc.status, apiErr.GetStatus())
			}

		})
	}
}
