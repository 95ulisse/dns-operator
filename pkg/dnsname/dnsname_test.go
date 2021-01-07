package dnsname

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	require := require.New(t)

	table := []struct {
		name    string
		success bool
	}{
		// No empty names
		{"", false},

		// Some good names
		{".", true},
		{"a", true},
		{"a.", true},
		{"a.b.c", true},
		{"a.b.c.", true},
		{"aa.bb.cc", true},
		{"a-a.b-b.c-c", true},

		// No double or leading dots
		{".a", false},
		{".a.b", false},
		{"..", false},
		{"a..", false},
		{"..a", false},
		{"a..b", false},
		{"a..b..", false},

		// Dashes at the beginning and end of a label are not allowed
		{"-a.b.c", false},
		{"a-.b.c", false},
		{"-a-.b.c", false},

		// Length tests
		{strings.Repeat("a", 63), true},
		{strings.Repeat("a", 64), false},
		{strings.Repeat("a.b.", 63) + "a", true},   // 253 chars
		{strings.Repeat("a.b.", 63) + "a.", true},  // 254 chars (with final dot)
		{strings.Repeat("a.b.", 63) + "ab", false}, // 254 chars (without final dot)
	}

	for _, entry := range table {
		name, err := NewName(entry.name)
		if entry.success {
			require.Nil(err, "Valid domain name %s failed parsing", entry.name)
			require.Equal(entry.name, name.String())
		} else {
			require.NotNil(err, "Invalid domain name %s has been correctly parsed", entry.name)
			require.IsType(&invalidNameError{}, err)
			require.Greater(len(err.Error()), 0)
		}
	}
}

func TestAccessors(t *testing.T) {
	require := require.New(t)

	// Root domain
	name, err := NewName(".")
	require.Nil(err)
	require.True(name.IsRoot())
	require.True(name.IsFQDN())

	// Simple FQDN
	name, err = NewName("example.com.")
	require.Nil(err)
	require.False(name.IsRoot())
	require.True(name.IsFQDN())

	// Simple domain
	name, err = NewName("example.com")
	require.Nil(err)
	require.False(name.IsRoot())
	require.False(name.IsFQDN())

	// FQDN conversion
	name = name.ToFQDN()
	require.Equal("example.com.", name.String())
	require.False(name.IsRoot())
	require.True(name.IsFQDN())

}

func TestIsChildOf(t *testing.T) {
	require := require.New(t)

	table := []struct {
		child    string
		parent   string
		expected bool
	}{
		{"example", "example.com", false},
		{"example.com", "example.com", true},
		{"example.com", "example2.com", false},
		{"example.com", "com", true},
		{"example.net", "com", false},
		{"example.net", ".", true},
		{"com", ".", true},
	}

	for _, entry := range table {

		// Parse child and parent
		child, err := NewName(entry.child)
		require.Nil(err)
		parent, err := NewName(entry.parent)
		require.Nil(err)

		// Test all possible combinations of FQDNs
		require.Equal(entry.expected, child.IsChildOf(parent), "Child: %s, Parent: %s", child.String(), parent.String())
		require.Equal(entry.expected, child.IsChildOf(parent.ToFQDN()), "Child: %s, Parent: %s", child.String(), parent.String())
		require.Equal(entry.expected, child.ToFQDN().IsChildOf(parent), "Child: %s, Parent: %s", child.String(), parent.String())
		require.Equal(entry.expected, child.ToFQDN().IsChildOf(parent.ToFQDN()), "Child: %s, Parent: %s", child.String(), parent.String())

	}

}

func TestJSON(t *testing.T) {
	require := require.New(t)

	// Test a valid name
	name, err := NewName("example.com")
	require.Nil(err)
	data, err := json.Marshal(&name)
	require.Nil(err)
	require.Equal("\"example.com\"", string(data))
	var name2 Name
	err = json.Unmarshal(data, &name2)
	require.Nil(err)
	require.Equal("example.com", name2.String())

	// Test an invalid name (we cannot build an invalid name, so we test only deserialization)
	var name3 Name
	err = json.Unmarshal([]byte("\"example..com\""), &name3)
	require.NotNil(err)

}
