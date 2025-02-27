/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCACFGAffilitionsApplyConfiguration represents an declarative configuration of the FabricCACFGAffilitions type for use
// with apply.
type FabricCACFGAffilitionsApplyConfiguration struct {
	AllowRemove *bool `json:"allowRemove,omitempty"`
}

// FabricCACFGAffilitionsApplyConfiguration constructs an declarative configuration of the FabricCACFGAffilitions type for use with
// apply.
func FabricCACFGAffilitions() *FabricCACFGAffilitionsApplyConfiguration {
	return &FabricCACFGAffilitionsApplyConfiguration{}
}

// WithAllowRemove sets the AllowRemove field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AllowRemove field is set to the value of the last call.
func (b *FabricCACFGAffilitionsApplyConfiguration) WithAllowRemove(value bool) *FabricCACFGAffilitionsApplyConfiguration {
	b.AllowRemove = &value
	return b
}
