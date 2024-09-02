package orberrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	msg := ErrInternalServerError.Error()
	expected := "500 Internal Server Error"
	require.Equal(t, expected, msg)
}

func TestWrappedError(t *testing.T) {
	err := ErrInternalServerError.Wrap(errors.New("testing"))
	expected := "500 Internal Server Error: testing"
	require.Equal(t, expected, err.Error())
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

func TestAs(t *testing.T) {
	orbe, ok := As(ErrRequestTimeout)
	require.True(t, ok)
	require.Equal(t, 408, orbe.Code)
}

func TestFromAndAs(t *testing.T) {
	err := From(errors.New("testing"))
	orbe, ok := As(err)
	require.True(t, ok)
	require.Equal(t, 500, orbe.Code)
}

func TestWrappedAs(t *testing.T) {
	err := ErrRequestTimeout.Wrap(errors.New("Test"))
	require.Equal(t, "408 Request Timeout: Test", err.Error())
	orbe, ok := As(err)
	require.True(t, ok)
	require.Equal(t, "408 Request Timeout: Test", orbe.Error())
	require.Equal(t, 408, orbe.Code)
}
