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
	JSON              bool
}

type Provider interface {
	BypassCF()
}