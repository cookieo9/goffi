go-ffi
======

Go FFI (and dlopen) packages to wrap C libraries.

This code is not being actively developped and is mostly provided as a reference.
For another FFI implementation that is active developpment see: https://bitbucket.org/binet/go-ffi

Automatically generated documentation can be found at:
 - http://go.pkgdoc.org/github.com/cookieo9/goffi/fcall   -- For the easy to use wrapper
 - http://go.pkgdoc.org/github.com/cookieo9/goffi/ffi     -- For lower level libffi access
 - http://go.pkgdoc.org/github.com/cookieo9/goffi/dlopen  -- For lower level dlopen access

Installation and Basic Usage
----------------------------

Install:

    go get github.com/cookieo9/goffi/fcall

Example (tested on Mac OS X):

    import "github.com/cookieo9/goffi/fcall"
    ...
    puts, _ := fcall.GetFunction("puts", SINT32, POINTER)
    cstr := fcall.CString("Hello, World!")
    defer fcall.Free(cstr)
    puts(cstr)

License
-------

http://cookieo9.mit-license.org/2012