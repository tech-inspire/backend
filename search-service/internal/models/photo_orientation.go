package models

type PhotoOrientation string

const (
	Portrait  PhotoOrientation = "portrait"
	Landscape PhotoOrientation = "landscape"
	Square    PhotoOrientation = "square"
)

func GetOrientationRange(orientation PhotoOrientation) (minRatio, maxRatio float64, ok bool) {
	const defaultTolerance = 0.05 // 5%

	switch orientation {
	case Square:
		return 1.0 - defaultTolerance, 1.0 + defaultTolerance, true
	case Portrait:
		return 0.0, 1.0 - defaultTolerance, true
	case Landscape:
		return 1.0 + defaultTolerance, 100.0, true // 100 as practical max
	default:
		return 0, 0, false
	}
}
