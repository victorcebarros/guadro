package guadro

import (
	"image"
)

type DisplayServer interface {
	Screenshot(geometry Geometry) (image.Image, error)
	MaxGeometry() (Geometry, error)
}
