package json

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadRejectsUnknownFields(t *testing.T) {
	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Keyboard","unexpected":true}`))
	var input struct {
		Name string `json:"name"`
	}

	if err := Read(request, &input); err == nil {
		t.Fatal("Read() error = nil, want an error for an unknown field")
	}
}
