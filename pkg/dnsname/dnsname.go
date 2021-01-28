package dnsname

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Use a regex to check the validity of labels and full domain names.
//
// Domain names are a sequence of dot-separated labels, each of which:
// - Can use only the following characters: [a-zA-Z0-9\.\-]
// - Does not start nor end with a dash
// - Is between 1 and 63 chars long
var nameRegexp = regexp.MustCompile("^(?:(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])\\.)*(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])\\.?$")

// Name represents a valid DNS resource name.
// +kubebuilder:validation:Type=string
type Name struct {
	// We have a json tag here only because kubebuilder cannot handle a struct without tags,
	// but we have our own custom (un)marshalling functions for this type.
	name string `json:",inline"`
}

// NewName validates the given domain name and constructs a new Name.
// The domain name can be either be fully qualified or not.
func NewName(name string) (*Name, error) {
	if !isValidDomainName(name) {
		return nil, &invalidNameError{name: name}
	}
	return &Name{name: name}, nil
}

// IsRoot returns `true` when the domain name is the root domain ".".
func (name *Name) IsRoot() bool {
	return name.name == "."
}

// IsFQDN return `true` when the domain name is a Fully Qualified Domain Name (i.e., it ends with a dot).
func (name *Name) IsFQDN() bool {
	return name.name[len(name.name)-1:] == "."
}

// ToFQDN makes this domain name fully qualified.
func (name *Name) ToFQDN() *Name {
	if name.IsFQDN() {
		return name
	}
	return &Name{name: name.name + "."}
}

// IsChildOf returns `true` if this domain name is a child of the given parent name.
// A domain is a child of another one if the latter is a suffix of the former.
// Note: this method ignores the final dot of a FQDN.
func (name *Name) IsChildOf(parent *Name) bool {
	childName := name.name
	if name.IsFQDN() {
		childName = childName[0 : len(childName)-1]
	}

	parentName := parent.name
	if parent.IsFQDN() {
		parentName = parentName[0 : len(parentName)-1]
	}

	return strings.HasSuffix(childName, parentName)
}

// String returns a string representation of this domain name.
func (name *Name) String() string {
	return name.name
}

// UnmarshalJSON parses the given JSON data as a domain name.
func (name *Name) UnmarshalJSON(data []byte) error {

	// Unmarshal the value as a string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Validate the domain name
	if !isValidDomainName(s) {
		return &invalidNameError{name: s}
	}

	name.name = s
	return nil

}

// MarshalJSON saves the JSON representation of this domain name.
func (name *Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(name.name)
}

func isValidDomainName(name string) bool {

	// No empty domain names
	if len(name) == 0 {
		return false
	}

	// Max domain length is 254 (including final "root" dot)
	if len(name) > 253 && (len(name) > 254 || !strings.HasSuffix(name, ".")) {
		return false
	}

	// As a special case, allow the root domain "."
	if name == "." {
		return true
	}

	// Validate the name using a regex
	if !nameRegexp.MatchString(name) {
		return false
	}

	return true

}
