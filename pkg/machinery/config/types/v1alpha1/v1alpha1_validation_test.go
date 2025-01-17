// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
)

type runtimeMode struct {
	requiresInstall bool
}

func (m runtimeMode) String() string {
	return fmt.Sprintf("runtimeMode(%v)", m.requiresInstall)
}

func (m runtimeMode) RequiresInstall() bool {
	return m.requiresInstall
}

func TestValidate(t *testing.T) {
	t.Parallel()

	endpointURL, err := url.Parse("https://localhost:6443/")
	require.NoError(t, err)

	for _, test := range []struct {
		name             string
		config           *v1alpha1.Config
		requiresInstall  bool
		strict           bool
		expectedWarnings []string
		expectedError    string
	}{
		{
			name: "NoMachine",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
			},
			expectedError: "1 error occurred:\n\t* machine instructions are required\n\n",
		},
		{
			name: "NoMachineType",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
				},
			},
			expectedWarnings: []string{
				`machine type is empty`,
			},
		},
		{
			name: "NoMachineTypeStrict",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
				},
			},
			strict:        true,
			expectedError: "1 error occurred:\n\t* warning: machine type is empty\n\n",
		},
		{
			name: "NoMachineInstall",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
				},
			},
		},
		{
			name: "NoMachineInstallRequired",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
				},
			},
			requiresInstall: true,
			expectedError:   "1 error occurred:\n\t* install instructions are required in \"runtimeMode(true)\" mode\n\n",
		},
		{
			name: "MachineInstallDisk",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
					MachineInstall: &v1alpha1.InstallConfig{
						InstallDisk: "/dev/vda",
					},
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
				},
			},
			requiresInstall: true,
		},

		{
			name: "ExternalCloudProviderEnabled",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
					ExternalCloudProviderConfig: &v1alpha1.ExternalCloudProviderConfig{
						ExternalEnabled: true,
						ExternalManifests: []string{
							"https://www.example.com/manifest1.yaml",
							"https://www.example.com/manifest2.yaml",
						},
					},
				},
			},
		},
		{
			name: "ExternalCloudProviderEnabledNoManifests",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
					ExternalCloudProviderConfig: &v1alpha1.ExternalCloudProviderConfig{
						ExternalEnabled: true,
					},
				},
			},
		},
		{
			name: "ExternalCloudProviderDisabled",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
					ExternalCloudProviderConfig: &v1alpha1.ExternalCloudProviderConfig{},
				},
			},
		},
		{
			name: "ExternalCloudProviderExtraManifests",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
					ExternalCloudProviderConfig: &v1alpha1.ExternalCloudProviderConfig{
						ExternalManifests: []string{
							"https://www.example.com/manifest1.yaml",
							"https://www.example.com/manifest2.yaml",
						},
					},
				},
			},
			expectedError: "1 error occurred:\n\t* external cloud provider is disabled, but manifests are provided\n\n",
		},
		{
			name: "ExternalCloudProviderInvalidManifests",
			config: &v1alpha1.Config{
				ConfigVersion: "v1alpha1",
				MachineConfig: &v1alpha1.MachineConfig{
					MachineType: "join",
				},
				ClusterConfig: &v1alpha1.ClusterConfig{
					ControlPlane: &v1alpha1.ControlPlaneConfig{
						Endpoint: &v1alpha1.Endpoint{
							endpointURL,
						},
					},
					ExternalCloudProviderConfig: &v1alpha1.ExternalCloudProviderConfig{
						ExternalEnabled: true,
						ExternalManifests: []string{
							"/manifest.yaml",
						},
					},
				},
			},
			expectedError: "1 error occurred:\n\t* invalid external cloud provider manifest url \"/manifest.yaml\": hostname must not be blank\n\n",
		},
	} {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			opts := []config.ValidationOption{config.WithLocal()}
			if test.strict {
				opts = append(opts, config.WithStrict())
			}

			warnings, errrors := test.config.Validate(runtimeMode{test.requiresInstall}, opts...)

			assert.Equal(t, test.expectedWarnings, warnings)

			if test.expectedError == "" {
				assert.NoError(t, errrors)
			} else {
				assert.EqualError(t, errrors, test.expectedError)
			}
		})
	}
}
