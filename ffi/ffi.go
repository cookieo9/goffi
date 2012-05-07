package ffi

// #cgo darwin CFLAGS: -I/usr/include/ffi
// #cgo LDFLAGS: -lffi
// #include <ffi.h>
import "C"

import (
	"unsafe"
	"reflect"
)

type Type *C.ffi_type

var (
	VOID       Type = &C.ffi_type_void
	UINT8      Type = &C.ffi_type_uint8
	SINT8      Type = &C.ffi_type_sint8
	UINT16     Type = &C.ffi_type_uint16
	SINT16     Type = &C.ffi_type_sint16
	UINT32     Type = &C.ffi_type_uint32
	SINT32     Type = &C.ffi_type_sint32
	UINT64     Type = &C.ffi_type_uint64
	SINT64     Type = &C.ffi_type_sint64
	FLOAT      Type = &C.ffi_type_float
	DOUBLE     Type = &C.ffi_type_double
	LONGDOUBLE Type = &C.ffi_type_longdouble
	POINTER    Type = &C.ffi_type_pointer
)

type ffi_error C.ffi_status

const (
	OK          = ffi_error(C.FFI_OK)
	BAD_TYPEDEF = ffi_error(C.FFI_BAD_TYPEDEF)
	BAD_ABI     = ffi_error(C.FFI_BAD_ABI)
)

func (fe ffi_error) Error() string {
	st := C.ffi_status(fe)
	switch st {
	case C.FFI_OK:
		return "NO ERROR!"
	case C.FFI_BAD_TYPEDEF:
		return "BAD_TYPEDEF"
	case C.FFI_BAD_ABI:
		return "BAD_ABI"
	}

	panic(fe)
}

type ABI C.ffi_abi

const (
	DEFAULT ABI = C.FFI_DEFAULT_ABI
)

type CIF C.ffi_cif

func (cif *CIF) Types () (rtype Type, atypes []Type) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&atypes))
	sh.Data = uintptr(unsafe.Pointer(cif.arg_types))
	sh.Cap = int(cif.nargs)
	sh.Len = int(cif.nargs)
	rtype = Type(cif.rtype)
	return
}


func NewCIF(abi ABI, rtype Type, atypes ...Type) (CIF, error) {
	var cif C.ffi_cif
	atypes_ptr := (**C.ffi_type)(unsafe.Pointer(&atypes[0]))
	status := C.ffi_prep_cif(&cif, C.ffi_abi(abi), C.uint(len(atypes)), rtype, atypes_ptr)
	if status == C.FFI_OK {
		return CIF(cif), nil
	}
	return CIF{}, ffi_error(status)
}

func Call(cif CIF, fn uintptr, rval uintptr, args ...uintptr) {
	C.ffi_call(
		(*C.ffi_cif)(&cif),
		(*[0]byte)(unsafe.Pointer(fn)),
		unsafe.Pointer(rval),
		(*unsafe.Pointer)(unsafe.Pointer(&args[0])),
	)
}
