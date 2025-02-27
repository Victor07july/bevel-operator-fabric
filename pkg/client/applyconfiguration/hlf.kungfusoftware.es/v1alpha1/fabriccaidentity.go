/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCAIdentityApplyConfiguration represents an declarative configuration of the FabricCAIdentity type for use
// with apply.
type FabricCAIdentityApplyConfiguration struct {
	Name        *string                                  `json:"name,omitempty"`
	Pass        *string                                  `json:"pass,omitempty"`
	Type        *string                                  `json:"type,omitempty"`
	Affiliation *string                                  `json:"affiliation,omitempty"`
	Attrs       *FabricCAIdentityAttrsApplyConfiguration `json:"attrs,omitempty"`
}

// FabricCAIdentityApplyConfiguration constructs an declarative configuration of the FabricCAIdentity type for use with
// apply.
func FabricCAIdentity() *FabricCAIdentityApplyConfiguration {
	return &FabricCAIdentityApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *FabricCAIdentityApplyConfiguration) WithName(value string) *FabricCAIdentityApplyConfiguration {
	b.Name = &value
	return b
}

// WithPass sets the Pass field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Pass field is set to the value of the last call.
func (b *FabricCAIdentityApplyConfiguration) WithPass(value string) *FabricCAIdentityApplyConfiguration {
	b.Pass = &value
	return b
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *FabricCAIdentityApplyConfiguration) WithType(value string) *FabricCAIdentityApplyConfiguration {
	b.Type = &value
	return b
}

// WithAffiliation sets the Affiliation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Affiliation field is set to the value of the last call.
func (b *FabricCAIdentityApplyConfiguration) WithAffiliation(value string) *FabricCAIdentityApplyConfiguration {
	b.Affiliation = &value
	return b
}

// WithAttrs sets the Attrs field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Attrs field is set to the value of the last call.
func (b *FabricCAIdentityApplyConfiguration) WithAttrs(value *FabricCAIdentityAttrsApplyConfiguration) *FabricCAIdentityApplyConfiguration {
	b.Attrs = value
	return b
}
