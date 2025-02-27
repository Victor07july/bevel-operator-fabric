/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricMainChannelExternalPeerOrganizationApplyConfiguration represents an declarative configuration of the FabricMainChannelExternalPeerOrganization type for use
// with apply.
type FabricMainChannelExternalPeerOrganizationApplyConfiguration struct {
	MSPID        *string `json:"mspID,omitempty"`
	TLSRootCert  *string `json:"tlsRootCert,omitempty"`
	SignRootCert *string `json:"signRootCert,omitempty"`
}

// FabricMainChannelExternalPeerOrganizationApplyConfiguration constructs an declarative configuration of the FabricMainChannelExternalPeerOrganization type for use with
// apply.
func FabricMainChannelExternalPeerOrganization() *FabricMainChannelExternalPeerOrganizationApplyConfiguration {
	return &FabricMainChannelExternalPeerOrganizationApplyConfiguration{}
}

// WithMSPID sets the MSPID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MSPID field is set to the value of the last call.
func (b *FabricMainChannelExternalPeerOrganizationApplyConfiguration) WithMSPID(value string) *FabricMainChannelExternalPeerOrganizationApplyConfiguration {
	b.MSPID = &value
	return b
}

// WithTLSRootCert sets the TLSRootCert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TLSRootCert field is set to the value of the last call.
func (b *FabricMainChannelExternalPeerOrganizationApplyConfiguration) WithTLSRootCert(value string) *FabricMainChannelExternalPeerOrganizationApplyConfiguration {
	b.TLSRootCert = &value
	return b
}

// WithSignRootCert sets the SignRootCert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SignRootCert field is set to the value of the last call.
func (b *FabricMainChannelExternalPeerOrganizationApplyConfiguration) WithSignRootCert(value string) *FabricMainChannelExternalPeerOrganizationApplyConfiguration {
	b.SignRootCert = &value
	return b
}
