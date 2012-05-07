package fcall

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// CString calls C.CString,
// so the user of this library doesn't have to
// import "C" / invoke cgo on their own code.
// The return value is a uintptr so the user doesn't
// have to deal with unsafe.Pointers.
func CString(s string) uintptr {
	return uintptr(unsafe.Pointer(C.CString(s)))
}

// Free calls C.free,
// so the user of this library doesn't have to
// import "C" / invoke cgo on their own code.
// The argument is a uintptr so the user doesn't
// have to deal with unsafe.Pointers.
func Free(ptr uintptr) {
	C.free(unsafe.Pointer(ptr))
}

// GoString calls C.GoString,
// so the user of this library doesn't have to
// import "C" / invoke cgo on their own code.
// The argument is a uintptr so the user doesn't
// have to deal with unsafe.Pointers.
func GoString(ptr uintptr) string {
	return C.GoString((*C.char)(unsafe.Pointer(ptr)))
}

// GoStringN calls C.GoStringN,
// so the user of this library doesn't have to
// import "C" / invoke cgo on their own code.
// The first argument is a uintptr so the user doesn't
// have to deal with unsafe.Pointers.
func GoStringN(ptr uintptr, n int32) string {
	return C.GoStringN((*C.char)(unsafe.Pointer(ptr)), C.int(n))
}
