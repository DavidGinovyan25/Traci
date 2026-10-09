package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAdminInput(t *testing.T) {
	var output bytes.Buffer
	input, err := readAdminInput(strings.NewReader(" admin \r\n Alice \n User \n ADMIN@example.com \n password123 \n"), &output)
	require.NoError(t, err)
	assert.Equal(t, "admin", input.Username)
	assert.Equal(t, "Alice", input.FirstName)
	assert.Equal(t, "User", input.SecondName)
	assert.Equal(t, "ADMIN@example.com", input.Email)
	assert.Equal(t, " password123 ", input.Password)
	assert.NotContains(t, output.String(), input.Password)
}

func TestReadAdminInputRejectsInvalidValues(t *testing.T) {
	for _, tc := range []struct{ name, input string }{
		{"missing input", ""},
		{"empty nickname", " \n"},
		{"long nickname", strings.Repeat("x", 51) + "\nA\nB\na@example.com\npassword123\n"},
		{"invalid email", "admin\nA\nB\nnot-email\npassword123\n"},
		{"display name email", "admin\nA\nB\nName <a@example.com>\npassword123\n"},
		{"short password", "admin\nA\nB\na@example.com\nshort\n"},
		{"long password", "admin\nA\nB\na@example.com\n" + strings.Repeat("x", 129) + "\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readAdminInput(strings.NewReader(tc.input), &bytes.Buffer{})
			require.Error(t, err)
		})
	}
}
