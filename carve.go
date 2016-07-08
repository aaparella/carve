// Package carve implements the seam carving algorithm for dyanmic image resizing.
package carve

import (
	"fmt"
	"image"
	"math"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

// ReduceHeight uses seam carving to reduce height of given image by n pixels.
func ReduceHeight(im image.Image, n int) (image.Image, error) {
	height := im.Bounds().Max.Y - im.Bounds().Min.Y
	if height < n {
		return im, fmt.Errorf("Cannot resize image of height %d by %d pixels", height, n)
	}

	for x := 0; x < n; x++ {
		energy := GenerateEnergyMap(im)
		seam := GenerateSeam(energy)
		im = RemoveSeam(im, seam)
	}
	return im, nil
}

// ReduceWidth uses seam carving to reduce width of given image by n pixels.
func ReduceWidth(im image.Image, n int) (image.Image, error) {
	width := im.Bounds().Max.Y - im.Bounds().Min.Y
	if width < n {
		return im, fmt.Errorf("Cannot resize image of width %d by %d pixels", width, n)
	}

	i := imaging.Rotate90(im)
	out, err := ReduceHeight(i, n)
	return imaging.Rotate270(out), err
}

// GenerateEnergyMap applies grayscale and sobel filters to the
// input image to create an energy map.
func GenerateEnergyMap(im image.Image) image.Image {
	g := gift.New(gift.Grayscale(), gift.Sobel())
	res := image.NewRGBA(im.Bounds())
	g.Draw(res, im)
	return res
}

// GenerateSeam returns the optimal horizontal seam for removal.
func GenerateSeam(im image.Image) Seam {
	mat := GenerateCostMatrix(im)
	return FindLowestCostSeam(mat)
}

// RemoveSeam creates a copy of the provided image, with the pixels at
// the points in the provided seam removed.
func RemoveSeam(im image.Image, seam Seam) image.Image {
	b := im.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()-1))
	min, max := b.Min, b.Max

	for _, point := range seam {
		x := point.X

		for y := min.Y; y < max.Y; y++ {
			if y == point.Y {
				continue
			}

			if y > point.Y {
				out.Set(x, y-1, im.At(x, y))
			} else {
				out.Set(x, y, im.At(x, y))
			}
		}
	}

	return out
}

// Seam defines a sequence of pixels through an image to be removed.
type Seam []Point

// Point defines an X Y point in an image.
type Point struct {
	X, Y int
}

// GenerateCostMatrix creates a matrix indicating the cumulative energy of the
// lowest cost seam from the left of the image to each pixel.
//
// mat[x][y] is the cumulative energy of the seam that runs from the left of
// the image to the pixel at column x, row y.
func GenerateCostMatrix(im image.Image) [][]float64 {
	min, max := im.Bounds().Min, im.Bounds().Max
	height, width := max.Y-min.Y, max.X-min.X

	mat := make([][]float64, width)
	for x := min.X; x < max.X; x++ {
		mat[x-min.X] = make([]float64, height)
	}

	for y := min.Y; y < max.Y; y++ {
		e, _, _, a := im.At(0, y).RGBA()
		mat[0][y-min.Y] = float64(e) / float64(a)
	}

	updatePoint := func(x, y int) {
		e, _, _, a := im.At(x, y).RGBA()

		up, down := math.MaxFloat64, math.MaxFloat64
		left := mat[x-1][y]
		if y != min.Y {
			up = mat[x-1][y-1]
		}
		if y < max.Y-1 {
			down = mat[x-1][y+1]
		}
		val := math.Min(float64(left), math.Min(float64(up), float64(down)))
		mat[x][y] = val + (float64(e) / float64(a))
	}

	// Calculate the remaining columns iteratively
	for x := min.X + 1; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			updatePoint(x, y)
		}
	}

	return mat
}

// FindLowestCostSeam uses a cumulative cost matrix to identify the seam with
// the lowest total cumulative energy.
func FindLowestCostSeam(mat [][]float64) Seam {
	width, height := len(mat), len(mat[0])
	seam := make([]Point, width)

	min, y := math.MaxFloat64, 0
	for ind, val := range mat[width-1] {
		if val < min {
			min = val
			y = ind
		}
	}

	seam[width-1] = Point{X: width - 1, Y: y}
	for x := width - 2; x >= 0; x-- {
		left := mat[x][y]
		up, down := math.MaxFloat64, math.MaxFloat64
		if y > 0 {
			up = mat[x][y-1]
		}
		if y < height-1 {
			down = mat[x][y+1]
		}

		if up <= left && up <= down {
			seam[x] = Point{X: x, Y: y - 1}
			y = y - 1
		} else if left <= up && left <= down {
			seam[x] = Point{X: x, Y: y}
		} else {
			seam[x] = Point{X: x, Y: y + 1}
			y = y + 1
		}
	}

	return seam
}
