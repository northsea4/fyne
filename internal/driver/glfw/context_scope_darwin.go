//go:build darwin

package glfw

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#include <stdint.h>
uintptr_t fyneBeginContextScope(void);
void fyneEndContextScope(uintptr_t scope);
*/
import "C"

func beginContextScope() uintptr    { return uintptr(C.fyneBeginContextScope()) }
func endContextScope(scope uintptr) { C.fyneEndContextScope(C.uintptr_t(scope)) }
