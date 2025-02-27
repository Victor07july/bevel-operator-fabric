/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

// FabricCASpecApplyConfiguration represents an declarative configuration of the FabricCASpec type for use
// with apply.
type FabricCASpecApplyConfiguration struct {
	Affinity         *v1.Affinity                           `json:"affinity,omitempty"`
	Tolerations      []v1.Toleration                        `json:"tolerations,omitempty"`
	ImagePullSecrets []v1.LocalObjectReference              `json:"imagePullSecrets,omitempty"`
	NodeSelector     *v1.NodeSelector                       `json:"nodeSelector,omitempty"`
	ServiceMonitor   *ServiceMonitorApplyConfiguration      `json:"serviceMonitor,omitempty"`
	Istio            *FabricIstioApplyConfiguration         `json:"istio,omitempty"`
	Database         *FabricCADatabaseApplyConfiguration    `json:"db,omitempty"`
	Hosts            []string                               `json:"hosts,omitempty"`
	Service          *FabricCASpecServiceApplyConfiguration `json:"service,omitempty"`
	Image            *string                                `json:"image,omitempty"`
	Version          *string                                `json:"version,omitempty"`
	Debug            *bool                                  `json:"debug,omitempty"`
	CLRSizeLimit     *int                                   `json:"clrSizeLimit,omitempty"`
	TLS              *FabricCATLSConfApplyConfiguration     `json:"rootCA,omitempty"`
	CA               *FabricCAItemConfApplyConfiguration    `json:"ca,omitempty"`
	TLSCA            *FabricCAItemConfApplyConfiguration    `json:"tlsCA,omitempty"`
	Cors             *CorsApplyConfiguration                `json:"cors,omitempty"`
	Resources        *v1.ResourceRequirements               `json:"resources,omitempty"`
	Storage          *StorageApplyConfiguration             `json:"storage,omitempty"`
	Metrics          *FabricCAMetricsApplyConfiguration     `json:"metrics,omitempty"`
	Env              []v1.EnvVar                            `json:"env,omitempty"`
}

// FabricCASpecApplyConfiguration constructs an declarative configuration of the FabricCASpec type for use with
// apply.
func FabricCASpec() *FabricCASpecApplyConfiguration {
	return &FabricCASpecApplyConfiguration{}
}

// WithAffinity sets the Affinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Affinity field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithAffinity(value v1.Affinity) *FabricCASpecApplyConfiguration {
	b.Affinity = &value
	return b
}

// WithTolerations adds the given value to the Tolerations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tolerations field.
func (b *FabricCASpecApplyConfiguration) WithTolerations(values ...v1.Toleration) *FabricCASpecApplyConfiguration {
	for i := range values {
		b.Tolerations = append(b.Tolerations, values[i])
	}
	return b
}

// WithImagePullSecrets adds the given value to the ImagePullSecrets field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ImagePullSecrets field.
func (b *FabricCASpecApplyConfiguration) WithImagePullSecrets(values ...v1.LocalObjectReference) *FabricCASpecApplyConfiguration {
	for i := range values {
		b.ImagePullSecrets = append(b.ImagePullSecrets, values[i])
	}
	return b
}

// WithNodeSelector sets the NodeSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NodeSelector field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithNodeSelector(value v1.NodeSelector) *FabricCASpecApplyConfiguration {
	b.NodeSelector = &value
	return b
}

// WithServiceMonitor sets the ServiceMonitor field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceMonitor field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithServiceMonitor(value *ServiceMonitorApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.ServiceMonitor = value
	return b
}

// WithIstio sets the Istio field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Istio field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithIstio(value *FabricIstioApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.Istio = value
	return b
}

// WithDatabase sets the Database field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Database field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithDatabase(value *FabricCADatabaseApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.Database = value
	return b
}

// WithHosts adds the given value to the Hosts field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Hosts field.
func (b *FabricCASpecApplyConfiguration) WithHosts(values ...string) *FabricCASpecApplyConfiguration {
	for i := range values {
		b.Hosts = append(b.Hosts, values[i])
	}
	return b
}

// WithService sets the Service field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Service field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithService(value *FabricCASpecServiceApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.Service = value
	return b
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithImage(value string) *FabricCASpecApplyConfiguration {
	b.Image = &value
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithVersion(value string) *FabricCASpecApplyConfiguration {
	b.Version = &value
	return b
}

// WithDebug sets the Debug field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Debug field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithDebug(value bool) *FabricCASpecApplyConfiguration {
	b.Debug = &value
	return b
}

// WithCLRSizeLimit sets the CLRSizeLimit field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CLRSizeLimit field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithCLRSizeLimit(value int) *FabricCASpecApplyConfiguration {
	b.CLRSizeLimit = &value
	return b
}

// WithTLS sets the TLS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TLS field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithTLS(value *FabricCATLSConfApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.TLS = value
	return b
}

// WithCA sets the CA field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CA field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithCA(value *FabricCAItemConfApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.CA = value
	return b
}

// WithTLSCA sets the TLSCA field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TLSCA field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithTLSCA(value *FabricCAItemConfApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.TLSCA = value
	return b
}

// WithCors sets the Cors field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Cors field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithCors(value *CorsApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.Cors = value
	return b
}

// WithResources sets the Resources field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Resources field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithResources(value v1.ResourceRequirements) *FabricCASpecApplyConfiguration {
	b.Resources = &value
	return b
}

// WithStorage sets the Storage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Storage field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithStorage(value *StorageApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.Storage = value
	return b
}

// WithMetrics sets the Metrics field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Metrics field is set to the value of the last call.
func (b *FabricCASpecApplyConfiguration) WithMetrics(value *FabricCAMetricsApplyConfiguration) *FabricCASpecApplyConfiguration {
	b.Metrics = value
	return b
}

// WithEnv adds the given value to the Env field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Env field.
func (b *FabricCASpecApplyConfiguration) WithEnv(values ...v1.EnvVar) *FabricCASpecApplyConfiguration {
	for i := range values {
		b.Env = append(b.Env, values[i])
	}
	return b
}
