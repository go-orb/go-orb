package orberrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	msg := ErrInternalServerError.Error()
	expected := "500 Internal Server Error"
	require.Equal(t, expected, msg)
}

func TestWrappedError(t *testing.T) {
	msg := fmt.Errorf("%w: %w", ErrInternalServerError, errors.New("testing")).Error()
	expected := "500 Internal Server Error: testing"
	require.Equal(t, expected, msg)
}

func TestNew(t *testing.T) {
	msg := New(500, "testing").Error()
	expected := "500 testing"
	require.Equal(t, expected, msg)
}

func TestNewHTTP(t *testing.T) {
	msg := NewHTTP(500).Error()
	expected := "500 Internal Server Error"
	require.Equal(t, expected, msg)
}

func TestFrom(t *testing.T) {
	msg := From(errors.New("testing")).Error()
	expected := "500 testing"
	require.Equal(t, expected, msg)
}
