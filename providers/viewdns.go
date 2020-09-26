package providers

type ViewDns struct {
	config *Config
}

func NewViewDns(c *Config) (Provider) {
	return &ViewDns{config: c}
}

func (v *ViewDns) BypassCF(domain string, results chan<- string) (error) {
	var err error
	return err
}