package fcall

import (
	"math"
	"runtime"
	. "testing"
)

func TestPuts(t *T) {
	puts, err := GetFunction("puts", SINT32, POINTER)
	if err != nil {
		t.Fatal(err)
	}

	cstr := CString("Hello, World!")
	defer Free(cstr)
	puts(cstr)
}

func TestSqrt(t *T) {
	var lib *Library
	switch runtime.GOOS {
	case "linux":
		mlib, err := OpenLibrary("libm.so", LAZY)
		if err != nil {
			t.Fatal(err)
		}
		lib = mlib
	case "darwin":
		lib = Default
	default:
		t.Fatal("Unknown System")
	}

	sqrt, err := lib.GetFunction("sqrt", DOUBLE, DOUBLE)
	if err != nil {
		t.Fatal(err)
	}

	x := sqrt(float64(9)).(float64)
	if math.Abs(x-3) > 0.0001 {
		t.Fatalf("Expected 3 got %g", x)
	}
	t.Logf("Got sqrt(9) = %g", x)
}
