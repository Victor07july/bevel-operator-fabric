/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCACFGApplyConfiguration represents an declarative configuration of the FabricCACFG type for use
// with apply.
type FabricCACFGApplyConfiguration struct {
	Identities   *FabricCACFGIdentitiesApplyConfiguration  `json:"identities,omitempty"`
	Affiliations *FabricCACFGAffilitionsApplyConfiguration `json:"affiliations,omitempty"`
}

// FabricCACFGApplyConfiguration constructs an declarative configuration of the FabricCACFG type for use with
// apply.
func FabricCACFG() *FabricCACFGApplyConfiguration {
	return &FabricCACFGApplyConfiguration{}
}

// WithIdentities sets the Identities field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Identities field is set to the value of the last call.
func (b *FabricCACFGApplyConfiguration) WithIdentities(value *FabricCACFGIdentitiesApplyConfiguration) *FabricCACFGApplyConfiguration {
	b.Identities = value
	return b
}

// WithAffiliations sets the Affiliations field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Affiliations field is set to the value of the last call.
func (b *FabricCACFGApplyConfiguration) WithAffiliations(value *FabricCACFGAffilitionsApplyConfiguration) *FabricCACFGApplyConfiguration {
	b.Affiliations = value
	return b
}
