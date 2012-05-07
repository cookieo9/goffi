package dlopen

// #cgo darwin CFLAGS: -D_DARWIN_C_SOURCE
// #cgo linux CFLAGS: -D_GNU_SOURCE
// #cgo linux LDFLAGS: -ldl
//
// #include <dlfcn.h>
// #include <stdlib.h>
//
// const void *rtld_default = RTLD_DEFAULT;
// const void *rtld_next = RTLD_NEXT;
import "C"

import (
	"errors"
	"unsafe"
)

type Flag C.int

const (
	LAZY Flag = C.RTLD_LAZY
	NOW  Flag = C.RTLD_NOW

	GLOBAL Flag = C.RTLD_GLOBAL
	LOCAL  Flag = C.RTLD_LOCAL
)

var (
	DEFAULT = uintptr(C.rtld_default)
	NEXT    = uintptr(C.rtld_next)
)

func Open(path string, flags Flag) (uintptr, error) {
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))

	lib := C.dlopen(cstr, C.int(flags))
	if lib == nil {
		err := dlerror("dlopen")
		return 0, err
	}

	return uintptr(lib), nil
}

func Symbol(handle uintptr, symbol string) (uintptr, error) {
	cstr := C.CString(symbol)
	defer C.free(unsafe.Pointer(cstr))

	sym := C.dlsym(unsafe.Pointer(handle), cstr)
	if sym == nil {
		return 0, dlerror("dlsym")
	}

	return uintptr(sym), nil
}

func dlerror(ctx string) error {
	errptr := C.dlerror()
	if errptr == nil {
		return nil
	}
	return errors.New(ctx + ": " + C.GoString(errptr))
}
