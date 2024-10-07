package objects

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/ulbios/bacnet/common"
)

func DecUnisgnedInteger(rawPayload APDUPayload) (uint32, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecUnisgnedInteger not ok: %v", rawPayload),
		)
	}

	if rawObject.TagNumber != TagUnsignedInteger || rawObject.TagClass {
		return 0, errors.Wrap(
			common.ErrWrongStructure,
			fmt.Sprintf("DecUnisgnedInteger wrong tag number: %v", rawObject.TagNumber),
		)
	}

	switch rawObject.Length {
	case 1:
		return uint32(rawObject.Data[0]), nil
	case 2:
		return uint32(binary.BigEndian.Uint16(rawObject.Data)), nil
	case 3:
		return uint32(uint16(uint32(rawObject.Data[0])<<16) | binary.BigEndian.Uint16(rawObject.Data[1:])), nil
	case 4:
		return binary.BigEndian.Uint32(rawObject.Data), nil
	}

	return 0, errors.Wrap(
		common.ErrNotImplemented,
		fmt.Sprintf("DecUnisgnedInteger not implemented data: %v", rawObject.Data),
	)
}
func EncUnsignedInteger8(value uint8) *Object {
	newObj := Object{}

	data := make([]byte, 1)
	data[0] = value

	newObj.TagNumber = TagUnsignedInteger
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}

func EncUnsignedInteger16(value uint16) *Object {
	newObj := Object{}

	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data[:], value)

	newObj.TagNumber = TagUnsignedInteger
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))
	
	return &newObj
}

func DecEnumerated(rawPayload APDUPayload) (uint32, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecEnumerated not ok %v", rawPayload),
		)
	}

	if rawObject.TagNumber != TagEnumerated || rawObject.TagClass {
		return 0, errors.Wrap(
			common.ErrWrongStructure,
			fmt.Sprintf("DecEnumerated wrong tag number: %v", rawObject.TagNumber),
		)
	}

	switch rawObject.Length {
	case 1:
		return uint32(rawObject.Data[0]), nil
	case 2:
		return uint32(binary.BigEndian.Uint16(rawObject.Data)), nil
	case 3:
		return uint32(uint16(uint32(rawObject.Data[0])<<16) | binary.BigEndian.Uint16(rawObject.Data[1:])), nil
	case 4:
		return binary.BigEndian.Uint32(rawObject.Data), nil
	}

	return 0, errors.Wrap(
		common.ErrNotImplemented,
		fmt.Sprintf("DecEnumerated not implemented data: %v", rawObject.Data),
	)
}

func EncEnumerated(value uint8) *Object {
	newObj := Object{}

	data := make([]byte, 1)
	data[0] = value

	newObj.TagNumber = TagEnumerated
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}

func DecReal(rawPayload APDUPayload) (float32, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecReal not ok: %v", rawPayload),
		)
	}

	if rawObject.TagNumber != TagReal {
		return 0, errors.Wrap(
			common.ErrWrongStructure,
			fmt.Sprintf("DecReal bad tag number: %v", rawObject.TagNumber),
		)
	}

	return math.Float32frombits(binary.BigEndian.Uint32(rawObject.Data)), nil
}

func EncReal(value float32) *Object {
	newObj := Object{}

	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data[:], math.Float32bits(value))

	newObj.TagNumber = TagReal
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}

func DecNull(rawPayload APDUPayload) (bool, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return false, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecNull not ok %v", rawPayload),
		)
	}

	if rawObject.TagNumber != TagReal {
		return false, errors.Wrap(
			common.ErrWrongStructure,
			fmt.Sprintf("DecNull bad tag number %v", rawObject.TagNumber),
		)
	}

	return rawObject.TagNumber == TagNull && !rawObject.TagClass && rawObject.Length == 0, nil
}

func EncNull() *Object {
	newObj := Object{}

	newObj.TagNumber = TagNull
	newObj.TagClass = false
	newObj.Data = nil
	newObj.Length = 0

	return &newObj
}
