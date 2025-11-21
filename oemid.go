// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type OEMID struct {
	// The type of oemid has multiple types, []byte or int
	//
	// oemid-type => oemid-pen / oemid-ieee / oemid-random
	// oemid-pen = int
	// oemid-ieee-cbor = bstr .size 3
	// oemid-ieee-json = base64-url-text .size 4
	// oemid-random-cbor = bstr .size 16
	// oemid-random-json = base64-url-text .size 24
	Value interface{}
}

func (o OEMID) Valid() bool {
	switch v := o.Value.(type) {
	case []byte:
		t := v
		return len(t) == 3 || len(t) == 16
	case int:
		return true
	default:
		return false
	}
}

func (o OEMID) MarshalJSON() ([]byte, error) {
	if !o.Valid() {
		return nil, fmt.Errorf("invalid value %#v", o.Value)
	}
	switch v := o.Value.(type) {
	case []byte:
		return json.Marshal(
			base64.RawURLEncoding.EncodeToString(v),
		)
	case int:
		return json.Marshal(o.Value)
	default:
		return nil, fmt.Errorf("invalid type %T", o.Value)
	}
}

func (o *OEMID) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		o.Value = i
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// must be base64url string
		value, err := base64.RawURLEncoding.DecodeString(s)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
		if len(value) != 3 && len(value) != 16 {
			return fmt.Errorf("%v is neither oemid-ieee nor oemid-random", data)
		}
		o.Value = value
		return nil
	}

	return fmt.Errorf("invalid value: expected oemid-pem, oemid-ieee or oemid-random")
}

func (o OEMID) MarshalCBOR() ([]byte, error) {
	if !o.Valid() {
		return nil, fmt.Errorf("invalid value %#v", o.Value)
	}
	switch o.Value.(type) {
	case []byte:
		return em.Marshal(o.Value)
	case int:
		return em.Marshal(o.Value)
	default:
		return nil, fmt.Errorf("invalid type %T", o.Value)
	}
}

func (o *OEMID) UnmarshalCBOR(data []byte) error {
	var i int
	if err := dm.Unmarshal(data, &i); err == nil {
		o.Value = i
		return nil
	}

	var b []byte
	if err := dm.Unmarshal(data, &b); err == nil {
		if len(b) != 3 && len(b) != 16 {
			return fmt.Errorf("%v is neither oemid-ieee nor oemid-random", data)
		}
		o.Value = b
		return nil
	}

	return fmt.Errorf("invalid value: expected oemid-pem, oemid-ieee or oemid-random")
}
