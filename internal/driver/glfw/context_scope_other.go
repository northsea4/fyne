//go:build !darwin

package glfw

func beginContextScope() uintptr { return 0 }
func endContextScope(uintptr)    {}
