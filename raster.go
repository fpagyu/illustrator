package illustrator

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/png"
	"log"
)

type Raster struct {
	Matrix            [6]float64
	Bounds            [4]float64
	Width             float64
	Height            float64
	Bits              int8 // bits per piexel in image map
	ImageType         int8 // 1=bitmap/grayscale; 3 = RGB; 4=CMYK
	AlphaChannelCount int8 // 0 = version 6.0; other values reserved for future versions
	BinAscii          int8 // 0=ASCII hexadecimal; 1 = binary
	ImageMask         int8 // 0 = opaque; 1 = transparent/colorized

	RawData []byte
}

func (r *Raster) B64Data() (string, error) {
	w, h := int(r.Width), int(r.Height)

	switch r.ImageType {
	case 1: // bitmap/grascale
		return r.Gray(w, h)
	case 3: // RGB
		return r.RGBA(w, h)
	case 4: // CMYK
		return r.CMYK(w, h)
	default:
		return "", errors.New("unknown image type")
	}
}

func (r *Raster) Gray(w, h int) (string, error) {
	header := "data:image/png:base64,"
	img := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := x + y*w
			img.Set(x, y, color.Gray{r.RawData[i]})
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		return "", err
	}

	return header + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (r *Raster) RGBA(w, h int) (string, error) {
	var header string
	var img = image.NewNRGBA(image.Rect(0, 0, w, h))
	if len(r.RawData)/(w*h) == 4 {
		alphaBase := w * h * 3
		header = "data:image/png;base64,"
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				off := x + y*w
				i := off * 3
				img.Set(x, y, color.RGBA{
					r.RawData[i], r.RawData[i+1],
					r.RawData[i+2], r.RawData[alphaBase+off],
				})
			}
		}
	} else {
		header = "data:image/jpeg;base64,"
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				i := (x + y*w) * 3
				img.Set(x, y, color.RGBA{
					r.RawData[i], r.RawData[i+1],
					r.RawData[i+2], 255,
				})
			}
		}
	}

	var buf = bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		return "", err
	}

	return header + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (r *Raster) CMYK(w, h int) (string, error) {
	var header = "data:image/png;base64,"
	img := image.NewCMYK(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (x + y*w) * 4
			img.Set(x, y, color.CMYK{
				r.RawData[i], r.RawData[i+1],
				r.RawData[i+2], r.RawData[i+3],
			})
		}
	}

	var buf = bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		return "", err
	}

	return header + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func XIArgs(vals []string) *Raster {
	if len(vals) < 20 {
		log.Println("invalid XI args:", vals)
		return nil
	}

	var obj Raster
	obj.Matrix[0] = toFloat(vals[1])
	obj.Matrix[1] = toFloat(vals[2])
	obj.Matrix[2] = toFloat(vals[3])
	obj.Matrix[3] = toFloat(vals[4])
	obj.Matrix[4] = toFloat(vals[5])
	obj.Matrix[5] = toFloat(vals[6])

	obj.Bounds[0] = toFloat(vals[8])
	obj.Bounds[1] = toFloat(vals[9])
	obj.Bounds[2] = toFloat(vals[10])
	obj.Bounds[3] = toFloat(vals[11])

	obj.Bits = toInt8(vals[14])
	obj.ImageType = toInt8(vals[15])
	obj.AlphaChannelCount = toInt8(vals[16])
	obj.BinAscii = toInt8(vals[18])
	obj.ImageMask = toInt8(vals[19])
	obj.Width, obj.Height = toFloat(vals[12]), toFloat(vals[13])

	return &obj
}
