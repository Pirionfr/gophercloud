package ssl

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// SSLTermConfig represents the SSL configuration for a particular load balancer.
type SSLTermConfig struct {
	// The port on which the SSL termination load balancer listens for secure
	// traffic. The value must be unique to the existing LB protocol/port
	// combination
	SecurePort int `mapstructure:"securePort"`

	// The private key for the SSL certificate which is validated and verified
	// against the provided certificates.
	PrivateKey string `mapstructure:"privatekey"`

	// The certificate used for SSL termination, which is validated and verified
	// against the key and intermediate certificate if provided.
	Certificate string

	// The intermediate certificate (for the user). The intermediate certificate
	// is validated and verified against the key and certificate credentials
	// provided. A user may only provide this value when accompanied by a
	// Certificate, PrivateKey, and SecurePort. It may not be added or updated as
	// a single attribute in a future operation.
	IntCertificate string `mapstructure:"intermediatecertificate"`

	// Determines if the load balancer is enabled to terminate SSL traffic or not.
	// If this is set to false, the load balancer retains its specified SSL
	// attributes but does not terminate SSL traffic.
	Enabled bool

	// Determines if the load balancer can only accept secure traffic. If set to
	// true, the load balancer will not accept non-secure traffic.
	SecureTrafficOnly bool
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	gophercloud.ErrResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult as a SSLTermConfig struct, if possible.
func (r GetResult) Extract() (*SSLTermConfig, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		SSL SSLTermConfig `mapstructure:"sslTermination"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.SSL, err
}

type CertificateMapping struct {
	ID             int
	HostName       string
	Certificate    string
	IntCertificate string `mapstructure:"intermediateCertificate"`
}

type CertMappingPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a CertMappingPage struct is empty.
func (p CertMappingPage) IsEmpty() (bool, error) {
	is, err := ExtractCertMappings(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractCertMappings accepts a Page struct, specifically a CertMappingPage struct, and extracts
// the elements into a slice of CertMapping structs. In other words, a generic
// collection is mapped into a relevant slice.
func ExtractCertMappings(page pagination.Page) ([]CertificateMapping, error) {
	type NestedMap struct {
		CertMap CertificateMapping `mapstructure:"certificateMapping" json:"certificateMapping"`
	}
	var resp struct {
		CertMappings []NestedMap `mapstructure:"certificateMappings" json:"certificateMappings"`
	}

	err := mapstructure.Decode(page.(CertMappingPage).Body, &resp)

	slice := []CertificateMapping{}

	for _, cert := range resp.CertMappings {
		slice = append(slice, cert.CertMap)
	}

	return slice, err
}

type CreateCertMappingResult struct {
	gophercloud.Result
}

// Extract interprets a result as a CertMapping struct, if possible.
func (r CreateCertMappingResult) Extract() (*CertificateMapping, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		CertMapping CertificateMapping `mapstructure:"certificateMapping"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.CertMapping, err
}
