package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
SetActivationPolicy(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
//
    return 0;
}

int
SetActivationPolicy2(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    return 0;
}
*/
import "C"
import "fmt"

func setActivationPolicy(turnOn bool) {
	fmt.Println("Setting ActivationPolicy")
	if turnOn {
		C.SetActivationPolicy2()
	} else {
		C.SetActivationPolicy()
	}

}
