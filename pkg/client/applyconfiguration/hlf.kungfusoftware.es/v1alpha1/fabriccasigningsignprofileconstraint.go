/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCASigningSignProfileConstraintApplyConfiguration represents an declarative configuration of the FabricCASigningSignProfileConstraint type for use
// with apply.
type FabricCASigningSignProfileConstraintApplyConfiguration struct {
	IsCA       *bool `json:"isCA,omitempty"`
	MaxPathLen *int  `json:"maxPathLen,omitempty"`
}

// FabricCASigningSignProfileConstraintApplyConfiguration constructs an declarative configuration of the FabricCASigningSignProfileConstraint type for use with
// apply.
func FabricCASigningSignProfileConstraint() *FabricCASigningSignProfileConstraintApplyConfiguration {
	return &FabricCASigningSignProfileConstraintApplyConfiguration{}
}

// WithIsCA sets the IsCA field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IsCA field is set to the value of the last call.
func (b *FabricCASigningSignProfileConstraintApplyConfiguration) WithIsCA(value bool) *FabricCASigningSignProfileConstraintApplyConfiguration {
	b.IsCA = &value
	return b
}

// WithMaxPathLen sets the MaxPathLen field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxPathLen field is set to the value of the last call.
func (b *FabricCASigningSignProfileConstraintApplyConfiguration) WithMaxPathLen(value int) *FabricCASigningSignProfileConstraintApplyConfiguration {
	b.MaxPathLen = &value
	return b
}
