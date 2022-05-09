// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MutableCbsd mutable cbsd
//
// swagger:model mutable_cbsd
type MutableCbsd struct {

	// capabilities
	// Required: true
	Capabilities Capabilities `json:"capabilities"`

	// desired state of cbsd in SAS
	// Required: true
	// Enum: [unregistered registered]
	DesiredState string `json:"desired_state"`

	// fcc id
	// Example: some_fcc_id
	// Required: true
	// Min Length: 1
	FccID string `json:"fcc_id"`

	// frequency preferences
	// Required: true
	FrequencyPreferences FrequencyPreferences `json:"frequency_preferences"`

	// serial number
	// Example: some_serial_number
	// Required: true
	// Min Length: 1
	SerialNumber string `json:"serial_number"`

	// user id
	// Example: some_user_id
	// Required: true
	// Min Length: 1
	UserID string `json:"user_id"`
}

// Validate validates this mutable cbsd
func (m *MutableCbsd) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCapabilities(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDesiredState(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFccID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFrequencyPreferences(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSerialNumber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MutableCbsd) validateCapabilities(formats strfmt.Registry) error {

	if err := m.Capabilities.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("capabilities")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("capabilities")
		}
		return err
	}

	return nil
}

var mutableCbsdTypeDesiredStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unregistered","registered"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		mutableCbsdTypeDesiredStatePropEnum = append(mutableCbsdTypeDesiredStatePropEnum, v)
	}
}

const (

	// MutableCbsdDesiredStateUnregistered captures enum value "unregistered"
	MutableCbsdDesiredStateUnregistered string = "unregistered"

	// MutableCbsdDesiredStateRegistered captures enum value "registered"
	MutableCbsdDesiredStateRegistered string = "registered"
)

// prop value enum
func (m *MutableCbsd) validateDesiredStateEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, mutableCbsdTypeDesiredStatePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *MutableCbsd) validateDesiredState(formats strfmt.Registry) error {

	if err := validate.RequiredString("desired_state", "body", m.DesiredState); err != nil {
		return err
	}

	// value enum
	if err := m.validateDesiredStateEnum("desired_state", "body", m.DesiredState); err != nil {
		return err
	}

	return nil
}

func (m *MutableCbsd) validateFccID(formats strfmt.Registry) error {

	if err := validate.RequiredString("fcc_id", "body", m.FccID); err != nil {
		return err
	}

	if err := validate.MinLength("fcc_id", "body", m.FccID, 1); err != nil {
		return err
	}

	return nil
}

func (m *MutableCbsd) validateFrequencyPreferences(formats strfmt.Registry) error {

	if err := m.FrequencyPreferences.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("frequency_preferences")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("frequency_preferences")
		}
		return err
	}

	return nil
}

func (m *MutableCbsd) validateSerialNumber(formats strfmt.Registry) error {

	if err := validate.RequiredString("serial_number", "body", m.SerialNumber); err != nil {
		return err
	}

	if err := validate.MinLength("serial_number", "body", m.SerialNumber, 1); err != nil {
		return err
	}

	return nil
}

func (m *MutableCbsd) validateUserID(formats strfmt.Registry) error {

	if err := validate.RequiredString("user_id", "body", m.UserID); err != nil {
		return err
	}

	if err := validate.MinLength("user_id", "body", m.UserID, 1); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this mutable cbsd based on the context it is used
func (m *MutableCbsd) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCapabilities(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateFrequencyPreferences(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MutableCbsd) contextValidateCapabilities(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Capabilities.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("capabilities")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("capabilities")
		}
		return err
	}

	return nil
}

func (m *MutableCbsd) contextValidateFrequencyPreferences(ctx context.Context, formats strfmt.Registry) error {

	if err := m.FrequencyPreferences.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("frequency_preferences")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("frequency_preferences")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MutableCbsd) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MutableCbsd) UnmarshalBinary(b []byte) error {
	var res MutableCbsd
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
