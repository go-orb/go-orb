package addr

import (
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYeet(t *testing.T) {
	t.Log(GetAddress(":0"))
}

func TestIPParser(t *testing.T) {
	var tests = []struct {
		IP       string
		Expected bool
	}{
		{":8080", false},
		{"192.168.1.1:8080", false},
		{"500.168.1.1:8080", true},
		{"[::]:8080", false},
		{"localhost:8080", false},
		{"", true},
		{"8080", true},
		{"192.168.1.1:808080808080", true},
		{"8080:", true},
		{"", true},
		{":abc", true},
		{"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080", false},
		{"[2001:db8::1]:8080", false},
	}

	for i, test := range tests {
		t.Run("TestIPParser"+strconv.Itoa(i), func(t *testing.T) {
			err := ValidateAddress(test.IP)
			assert.Equal(t, test.Expected, err != nil, test.IP, err)
		})
	}
}

func TestParsePort(t *testing.T) {
	_, err := ParsePort(":8080")
	assert.NoError(t, err)
	_, err = ParsePort(":abc")
	assert.Error(t, err)
}

func TestIsLocal(t *testing.T) {
	testData := []struct {
		addr   string
		expect bool
	}{
		{"localhost", true},
		{"localhost:8080", true},
		{"127.0.0.1", true},
		{"127.0.0.1:1001", true},
		{"80.1.1.1", false},
	}

	for _, d := range testData {
		res := IsLocal(d.addr)
		if res != d.expect {
			t.Fatalf("expected %t got %t", d.expect, res)
		}
	}
}

func TestExtractor(t *testing.T) {
	testData := []struct {
		addr   string
		expect string
		parse  bool
	}{
		{"127.0.0.1", "127.0.0.1", false},
		{"10.0.0.1", "10.0.0.1", false},
		{"", "", true},
		{"0.0.0.0", "", true},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "", true},
	}

	for _, d := range testData {
		addr, err := Extract(d.addr)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if d.parse {
			ip := net.ParseIP(addr)
			if ip == nil {
				t.Error("Unexpected nil IP for " + addr)
			}
		} else if addr != d.expect {
			t.Errorf("Expected %s got %s", d.expect, addr)
		}
	}
}
