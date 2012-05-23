// Package fcall wraps the ffi and dlopen packages to provide easy access to foreign function calls from Go Code.
package fcall

import (
	"fmt"
	"github.com/cookieo9/goffi/dlopen"
	"github.com/cookieo9/goffi/ffi"
	"reflect"
)

// Flags passed to OpenLibrary.
type OpenFlag dlopen.Flag

var (
	LAZY = OpenFlag(dlopen.LAZY) // Bind Symbols on Use
	NOW  = OpenFlag(dlopen.NOW)  // Bind Symbols on Open

	GLOBAL = OpenFlag(dlopen.GLOBAL) // Load Symbols into Global Namespace
	LOCAL  = OpenFlag(dlopen.LOCAL)  // Load Symbols into Private Namespace
)

// Types for parameters and return values to/from
// FFI called functions.
type Type ffi.Type

var (
	VOID    = Type(ffi.VOID)
	UINT8   = Type(ffi.UINT8)
	UINT16  = Type(ffi.UINT16)
	UINT32  = Type(ffi.UINT32)
	UINT64  = Type(ffi.UINT64)
	SINT8   = Type(ffi.SINT8)
	SINT16  = Type(ffi.SINT16)
	SINT32  = Type(ffi.SINT32)
	SINT64  = Type(ffi.SINT64)
	FLOAT   = Type(ffi.FLOAT)
	DOUBLE  = Type(ffi.DOUBLE)
	POINTER = Type(ffi.POINTER)
)

func convert_type(t Type) reflect.Type {
	switch t {
	case VOID:
		return reflect.TypeOf(nil)
	case UINT8:
		return reflect.TypeOf(uint8(0))
	case UINT16:
		return reflect.TypeOf(uint16(0))
	case UINT32:
		return reflect.TypeOf(uint32(0))
	case UINT64:
		return reflect.TypeOf(uint64(0))
	case SINT8:
		return reflect.TypeOf(int8(0))
	case SINT16:
		return reflect.TypeOf(int16(0))
	case SINT32:
		return reflect.TypeOf(int32(0))
	case SINT64:
		return reflect.TypeOf(int64(0))
	case FLOAT:
		return reflect.TypeOf(float32(0))
	case DOUBLE:
		return reflect.TypeOf(float64(0))
	case POINTER:
		return reflect.TypeOf(uintptr(0))
	}
	return nil
}

// A Library represents an dlopen'd library from
// which functions can be loaded & wrapped.
type Library struct {
	lib uintptr
}

var (
	Default = &Library{dlopen.DEFAULT} // Default Library (RTLD_DEFAULT)
	Next    = &Library{dlopen.NEXT}    // Next Library (RTLD_NEXT)
)

// OpenLibrary opens a shared-library / DLL / shared object
// using dlopen.
func OpenLibrary(path string, flags OpenFlag) (*Library, error) {
	lib, err := dlopen.Open(path, dlopen.Flag(flags))
	if err != nil {
		return nil, err
	}

	return &Library{lib: lib}, nil
}

// A Function is a closure that represents a wrapped C function.
// Calling the closure checks the parameters against those registered
// into the wrapper, and will panic if any mismatches occur.
type Function func(args ...interface{}) interface{}

// GetFunction returns a Go function (produced by Wrap()) that
// can dynamically call the named function in the current library.
func (lib *Library) GetFunction(symbol string, rtype Type, atypes ...Type) (Function, error) {
	sym, err := dlopen.Symbol(lib.lib, symbol)
	if err != nil {
		return nil, err
	}
	return Wrap(sym, rtype, atypes...)
}

// Wrap creates a callable function from a C function pointer
// and a list of parameter and return types. This method is 
// exported for users that wish to wrap a C function pointer
// aquired from the raw dlopen library, or through other means.
// 
// The default calling convention for the platform (as defined by ffi)
// is used for the call.
func Wrap(fn uintptr, rtype Type, atypes ...Type) (Function, error) {
	ffiatypes := make([]ffi.Type, len(atypes))
	for i, t := range atypes {
		ffiatypes[i] = ffi.Type(t)
	}

	cif, err := ffi.NewCIF(ffi.DEFAULT, ffi.Type(rtype), ffiatypes...)
	if err != nil {
		return nil, err
	}
	return WrapCIF(fn, cif)
}

// WrapCIF creates a callable function from a C function pointer
// and a ffi.CIF structure, thus allowing a user to call a
// function using one of libffi's other supported calling
// conventions.
func WrapCIF(fn uintptr, cif ffi.CIF) (Function, error) {
	rtype0, atypes0 := cif.Types()

	rtype := convert_type(Type(rtype0))
	atypes := make([]reflect.Type, len(atypes0))
	for i, atype := range atypes0 {
		if t := convert_type(Type(atype)); t != nil {
			atypes[i] = t
		} else {
			return nil, fmt.Errorf("ffi_cif: unsupported type for arg %d", i+1)
		}
	}

	f := func(args ...interface{}) interface{} {
		if len(args) != len(atypes) {
			panic(fmt.Errorf("Expecting %d args, got %d instead", len(atypes), len(args)))
		}

		aptrs := make([]uintptr, len(args))

		for i, arg := range args {
			aval := reflect.ValueOf(arg)
			if aval.Type() != atypes[i] {
				panic(fmt.Errorf("Expecting arg %d to be %T but got %T instead", i+1, atypes[i], aval.Type()))
			}
			aval2 := reflect.New(aval.Type())
			aval2.Elem().Set(aval)
			aptrs[i] = aval2.Pointer()
		}

		if Type(rtype0) == VOID {
			ffi.Call(cif, fn, 0, aptrs...)
			return nil
		}

		rval := reflect.New(rtype)
		ffi.Call(cif, fn, rval.Pointer(), aptrs...)
		return rval.Elem().Interface()
	}
	return f, nil
}

// Calls Default.GetFunction() to create a wrapped function for a symbol in
// the default library.
func GetFunction(name string, rtype Type, atypes ...Type) (Function, error) {
	return Default.GetFunction(name, rtype, atypes...)
}
