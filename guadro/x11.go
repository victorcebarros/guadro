package guadro

// #cgo LDFLAGS: -lX11
// #include <stdlib.h>
// #include <X11/Xlib.h>
// #include <X11/Xutil.h>
//
// int goXDestroyImage(XImage *ximage) {
//     return XDestroyImage(ximage); // may be defined as macro
// }
//
// unsigned long goXGetPixel(XImage *ximage, int x, int y) {
//     return XGetPixel(ximage, x, y); // may be defined as macro
// }
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type X11Server struct {
	Display string
}

func (x11 X11Server) Screenshot(geometry Geometry) (image.Image, error) {

	displayName := C.CString(x11.Display)
	defer C.free(unsafe.Pointer(displayName))

	display := C.XOpenDisplay(displayName)

	if display == nil {
		display := x11.Display

		if display == "" {
			display = ":0"
		}

		return nil, fmt.Errorf("Could not open display %s", display)
	}

	defer C.XCloseDisplay(display)

	window := C.XDefaultRootWindow(display)

	// defined in X11/X.h as a macro
	ZPixmap := C.int(2)

	xImage := C.XGetImage(
		display, window,
		C.int(geometry.XOffset), C.int(geometry.YOffset),
		C.uint(geometry.Width), C.uint(geometry.Height),
		C.XAllPlanes(), ZPixmap,
	)

	defer C.goXDestroyImage(xImage)

	RMask := xImage.red_mask
	GMask := xImage.green_mask
	BMask := xImage.blue_mask

	rect := image.Rect(0, 0, geometry.Width, geometry.Height)

	screenshot := image.NewRGBA(rect)

	for x := 0; x < geometry.Width; x++ {
		for y := 0; y < geometry.Height; y++ {
			xPixel := C.goXGetPixel(xImage, C.int(x), C.int(y))

			pixel := color.RGBA{
				R: uint8(uint64(xPixel&RMask) >> 16),
				G: uint8(uint64(xPixel&GMask) >> 8),
				B: uint8(xPixel & BMask),
				A: (1 << 8) - 1,
			}

			screenshot.Set(x, y, pixel)
		}
	}

	return screenshot, nil
}

func (x11 X11Server) MaxGeometry() (Geometry, error) {
	displayName := C.CString(x11.Display)
	defer C.free(unsafe.Pointer(displayName))

	display := C.XOpenDisplay(displayName)

	if display == nil {
		return Geometry{}, fmt.Errorf("Could not open display %s", x11.Display)
	}

	defer C.XCloseDisplay(display)

	window := C.XDefaultRootWindow(display)

	var xattrs C.XWindowAttributes

	// Yes, zero is the value returned in case of failure
	if C.XGetWindowAttributes(display, window, &xattrs) == 0 {
		return Geometry{}, fmt.Errorf("Failed to get Window Attributes")
	}

	geometry := Geometry{
		Width:  int(xattrs.width),
		Height: int(xattrs.height),
	}

	return geometry, nil
}
