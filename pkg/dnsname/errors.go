package dnsname

import "fmt"

type invalidNameError struct {
	name string
}

func (e *invalidNameError) Error() string {
	if len(e.name) == 0 {
		return "Invalid empty domain name"
	}
	return fmt.Sprintf("Invalid domain name: %s", e.name)
}
