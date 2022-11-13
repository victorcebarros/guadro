package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/victorcebarros/guadro/guadro"
)

func init() {
	pflag.ErrHelp = nil
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [OPTIONs]...\n\n", os.Args[0])
		pflag.PrintDefaults()
		fmt.Println()
		fmt.Fprintln(os.Stderr,
			"geometry is a standard X string, see man 3 XParseGeometry")
		fmt.Fprintln(os.Stderr,
			"format is one of the following: png")
		// TODO: implement support for jpeg, bmp, tiff, webp formats
		os.Exit(1)
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "guadro: error: %s\n", err)
	os.Exit(1)
}

// TODO: remove responsibility from main function
func main() {
	geometryStr := pflag.StringP("geometry", "g", "", "sets geometry of area to screenshot")
	outfile := pflag.StringP("output", "o", "", "sets output file name")
	format := pflag.StringP("format", "f", "png", "sets output format to file")
	pflag.Parse()

	// TODO: untangle guadro from X11
	// TODO: pass DisplayServer specific flags so it can parse itself
	x11 := guadro.X11Server{}
	geometry, err := x11.MaxGeometry()

	if err != nil {
		die(err)
	}

	if *geometryStr != "" {
		geometry, err = guadro.ParseGeometry(*geometryStr)

		if err != nil {
			die(err)
		}
	}

	if *outfile == "" {
		*outfile = fmt.Sprintf("%d.%s", time.Now().Unix(), *format)
	}

	img, err := x11.Screenshot(geometry)

	if err != nil {
		die(err)
	}

	file, err := os.OpenFile(*outfile, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		die(err)
	}

	defer file.Close()

	err = png.Encode(file, img)

	if err != nil {
		die(err)
	}
}
