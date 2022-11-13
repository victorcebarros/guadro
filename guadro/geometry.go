package guadro

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Geometry struct {
	Width, Height    int
	XOffset, YOffset int
}

func parseNumberWithSign(s string) (int, int, error) {
	if s == "" {
		return 0, 0, errors.New("missing offset")
	}

	sign := 0

	switch s[0] {
	case '+':
		sign = +1
	case '-':
		sign = -1
	default:
		return 0, 0, errors.New("invalid sign")
	}

	s = s[1:]

	if s[0] == '+' || s[0] == '-' {
		return 0, 0, errors.New("multiple signs")
	}

	number := 0

	_, err := fmt.Sscanf(s, "%d", &number)

	if err != nil {
		return 0, 0, err
	}

	// + 1 for the sign
	bytes := len(strconv.Itoa(number)) + 1

	number *= sign

	return number, bytes, nil
}

func parseXYOffsets(geometry *Geometry, s string) error {
	if s == "" {
		return nil
	}

	n, b, err := parseNumberWithSign(s)

	if err != nil {
		return err
	}

	geometry.XOffset = n

	s = s[b:]

	n, b, err = parseNumberWithSign(s)

	if err != nil {
		return err
	}

	geometry.YOffset = n

	s = s[b:]

	if s != "" {
		return errors.New("trailing garbage on geometry string")
	}

	return nil
}

func parseWidthAndHeight(geometry *Geometry, s string) error {
	_, err := fmt.Sscanf(s, "%d", &geometry.Width)

	if err != nil {
		return err
	}

	// we consider that there will never be any signs in
	// the beginning of the string
	n := len(strconv.Itoa(geometry.Width))

	s = s[n:]

	if s == "" {
		return errors.New("incomplete width and heigth")
	}

	if s[0] != 'x' && s[0] != 'X' {
		return errors.New("missing {xX} separator")
	}

	s = s[1:]

	if s == "" {
		return errors.New("missing height")
	}

	if !unicode.IsDigit(rune(s[0])) {
		return errors.New("unexpected sign on height")
	}

	_, err = fmt.Sscanf(s, "%d", &geometry.Height)

	if err != nil {
		return err
	}

	n = len(strconv.Itoa(geometry.Height))

	return parseXYOffsets(geometry, s[n:])
}

func parseGeometry(geometry *Geometry, s string) error {
	// we consider that both "" and "=" are valid ss
	if s == "" {
		return nil
	}

	if s[0] == '=' {
		s = s[1:]
	}

	if s == "" {
		return nil
	}

	if s[0] == '+' || s[0] == '-' {
		// can only contain XY offsets
		return parseXYOffsets(geometry, s)
	}

	// will contain Width and Height, and may contain XY offsets
	return parseWidthAndHeight(geometry, s)
}

// Geometry size and placement follows X's standard strings.
// Refer to $ man 3 XParseGeometry

// [=][<width>{xX}<height>][{+-}<xoffset>{+-}<yoffset>]
//     [.*] => Optional
//     {.*} => Either of
//     <.*> => Integer
// For example:
//    // ignoring errors
//    g, _ := ParseGeometry("=200x300+103-44")
//    // prints guadro.Geometry{Width:200, Height:300, XOffset:103, YOffset:-44}
//    fmt.Printf("%#v\n", g)

func ParseGeometry(s string) (Geometry, error) {
	geometry := Geometry{}
	err := parseGeometry(&geometry, s)
	return geometry, err
}
