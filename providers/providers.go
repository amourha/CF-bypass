package providers

import (
	"io"
	"../http"
)

type Config struct {
	Verbose           bool
	Client            *http.Client
	Providers         []string
	Output            io.Writer
}

type Provider interface {
	BypassCF(domain string, results chan<- string) error
}