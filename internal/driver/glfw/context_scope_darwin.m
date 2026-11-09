#import <Foundation/Foundation.h>
#include <stdint.h>

// The retained pool spans synchronous Go GL calls on the locked GLFW thread.
// A local @autoreleasepool would drain before those calls begin. RunWithContext
// pairs this owned token with a defer, including when the Go callback panics,
// and drains it before detaching the still-live GL context.
uintptr_t fyneBeginContextScope(void) {
    return (uintptr_t)[[NSAutoreleasePool alloc] init];
}

void fyneEndContextScope(uintptr_t scope) {
    [(NSAutoreleasePool *)scope drain];
}
