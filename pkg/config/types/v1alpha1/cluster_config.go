/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha1

import (
	"net/url"
	"strings"

	"github.com/talos-systems/talos/pkg/config/cluster"
	"github.com/talos-systems/talos/pkg/config/machine"
	"github.com/talos-systems/talos/pkg/constants"
	"github.com/talos-systems/talos/pkg/crypto/x509"
)

// ClusterConfig reperesents the cluster-wide config values
type ClusterConfig struct {
	ControlPlane                  *ControlPlaneConfig               `yaml:"controlPlane"`
	ClusterName                   string                            `yaml:"clusterName,omitempty"`
	ClusterNetwork                *ClusterNetworkConfig             `yaml:"network,omitempty"`
	BootstrapToken                string                            `yaml:"token,omitempty"`
	CertificateKey                string                            `yaml:"certificateKey"`
	ClusterAESCBCEncryptionSecret string                            `yaml:"aescbcEncryptionSecret"`
	ClusterCA                     *x509.PEMEncodedCertificateAndKey `yaml:"ca,omitempty"`
	APIServer                     *APIServerConfig                  `yaml:"apiServer,omitempty"`
	ControllerManager             *ControllerManagerConfig          `yaml:"controllerManager,omitempty"`
	Scheduler                     *SchedulerConfig                  `yaml:"scheduler,omitempty"`
	EtcdConfig                    *EtcdConfig                       `yaml:"etcd,omitempty"`
}

// Endpoint struct holds the endpoint url parsed out of machine config
type Endpoint struct {
	*url.URL
}

// UnmarshalYAML is a custom unmarshaller for the endpoint struct
func (e *Endpoint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var endpoint string

	if err := unmarshal(&endpoint); err != nil {
		return err
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	*e = Endpoint{url}

	return nil
}

// MarshalYAML is a custom unmarshaller for the endpoint struct
func (e *Endpoint) MarshalYAML() (interface{}, error) {
	return e.URL.String(), nil
}

// ControlPlaneConfig represents control plane config vals
type ControlPlaneConfig struct {
	Version string `yaml:"version"`

	// Endpoint is the canonical controlplane endpoint, which can be an IP
	// address or a DNS hostname, is single-valued, and may optionally include a
	// port number.
	Endpoint *Endpoint `yaml:"endpoint"`

	// LocalAPIServerPort is the port that the api server listens to internally.
	// This may be different than the port portion listed in the endpoint field above.
	LocalAPIServerPort int `yaml:"localAPIServerPort,omitempty"`
}

// APIServerConfig represents kube apiserver config vals
type APIServerConfig struct {
	Image     string            `yaml:"image,omitempty"`
	ExtraArgs map[string]string `yaml:"extraArgs,omitempty"`
	CertSANs  []string          `yaml:"certSANs,omitempty"`
}

// ControllerManagerConfig represents kube controller manager config vals
type ControllerManagerConfig struct {
	Image     string            `yaml:"image,omitempty"`
	ExtraArgs map[string]string `yaml:"extraArgs,omitempty"`
}

// SchedulerConfig represents kube scheduler config vals
type SchedulerConfig struct {
	Image     string            `yaml:"image,omitempty"`
	ExtraArgs map[string]string `yaml:"extraArgs,omitempty"`
}

// EtcdConfig represents etcd config vals
type EtcdConfig struct {
	ContainerImage string                            `yaml:"image,omitempty"`
	RootCA         *x509.PEMEncodedCertificateAndKey `yaml:"ca"`
}

// ClusterNetworkConfig represents kube networking config vals
type ClusterNetworkConfig struct {
	CNI           string   `yaml:"cni"`
	DNSDomain     string   `yaml:"dnsDomain"`
	PodSubnet     []string `yaml:"podSubnets"`
	ServiceSubnet []string `yaml:"serviceSubnets"`
}

// Version implements the Configurator interface.
func (c *ClusterConfig) Version() string {
	return c.ControlPlane.Version
}

// Endpoint implements the Configurator interface.
func (c *ClusterConfig) Endpoint() *url.URL {
	return c.ControlPlane.Endpoint.URL
}

// LocalAPIServerPort implements the Configurator interface.
func (c *ClusterConfig) LocalAPIServerPort() int {
	if c.ControlPlane.LocalAPIServerPort == 0 {
		return 6443
	}

	return c.ControlPlane.LocalAPIServerPort
}

// CertSANs implements the Configurator interface.
func (c *ClusterConfig) CertSANs() []string {
	return c.APIServer.CertSANs
}

// SetCertSANs implements the Configurator interface.
func (c *ClusterConfig) SetCertSANs(sans []string) {
	if c.APIServer == nil {
		c.APIServer = &APIServerConfig{}
	}

	c.APIServer.CertSANs = append(c.APIServer.CertSANs, sans...)
}

// CA implements the Configurator interface.
func (c *ClusterConfig) CA() *x509.PEMEncodedCertificateAndKey {
	return c.ClusterCA
}

// AESCBCEncryptionSecret implements the Configurator interface.
func (c *ClusterConfig) AESCBCEncryptionSecret() string {
	return c.ClusterAESCBCEncryptionSecret
}

// Config implements the Configurator interface.
func (c *ClusterConfig) Config(t machine.Type) (string, error) {
	return "", nil
}

// Etcd implements the Configurator interface.
func (c *ClusterConfig) Etcd() cluster.Etcd {
	return c.EtcdConfig
}

// Image implements the Configurator interface.
func (e *EtcdConfig) Image() string {
	return e.ContainerImage
}

// CA implements the Configurator interface.
func (e *EtcdConfig) CA() *x509.PEMEncodedCertificateAndKey {
	return e.RootCA
}

// Token implements the Configurator interface.
func (c *ClusterConfig) Token() cluster.Token {
	return c
}

// ID implements the Configurator interface.
func (c *ClusterConfig) ID() string {
	parts := strings.Split(c.BootstrapToken, ".")
	if len(parts) != 2 {
		return ""
	}

	return parts[0]
}

// Secret implements the Configurator interface.
func (c *ClusterConfig) Secret() string {
	parts := strings.Split(c.BootstrapToken, ".")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

// Network implements the Configurator interface.
func (c *ClusterConfig) Network() cluster.Network {
	return c
}

// CNI implements the Configurator interface.
func (c *ClusterConfig) CNI() string {
	switch {
	case c.ClusterNetwork == nil:
		fallthrough
	case c.ClusterNetwork.CNI == "":
		return constants.DefaultCNI
	}

	return c.ClusterNetwork.CNI
}

// PodCIDR implements the Configurator interface.
func (c *ClusterConfig) PodCIDR() string {
	switch {
	case c.ClusterNetwork == nil:
		fallthrough
	case len(c.ClusterNetwork.PodSubnet) == 0:
		return constants.DefaultPodCIDR
	}

	return c.ClusterNetwork.PodSubnet[0]
}

// ServiceCIDR implements the Configurator interface.
func (c *ClusterConfig) ServiceCIDR() string {
	switch {
	case c.ClusterNetwork == nil:
		fallthrough
	case len(c.ClusterNetwork.ServiceSubnet) == 0:
		return constants.DefaultServiceCIDR
	}

	return c.ClusterNetwork.ServiceSubnet[0]
}