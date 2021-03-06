package imagick

import (
	. "github.com/journeymidnight/pipa/error"
	"github.com/journeymidnight/pipa/helper"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type ResizePlan struct {
	Mode                string
	Width               int
	Height              int
	Long                int
	Short               int
	Limit               bool
	Color               string
	WatermarkProportion int
	Proportion          int
	Data                []byte
}

type WatermarkPlan struct {
	Transparency  int
	RotateDegrees int
	Position      string
	XMargin       int
	YMargin       int
	Voffset       int
	PictureMask   WatermarkPicture
	TextMask      WatermarkText
	Order         int //default 0:Image watermark in front 1:Text watermark in front
	Align         int //default 0:align from the top 1:middle 2:bottom
	Interval      int
}

type WatermarkPicture struct {
	Bucket string
	Image  string
	Data   []byte
	Crop   CropPlan
}

type WatermarkText struct {
	Text   string
	Type   string
	Color  string
	Size   int
	Shadow int
	Fill   bool
}

type CropPlan struct {
}

type RotatePlan struct {
	Degrees int
	Color   string
}

type ImageWand struct {
	MagickWand *imagick.MagickWand
	PixelWand  *imagick.PixelWand
	DrawWand   *imagick.DrawingWand
}

func Initialize() (lib interface{}) {
	imagick.Initialize()
	return lib
}

func (img *ImageWand) Destory() {
	img.MagickWand.Destroy()
	img.DrawWand.Destroy()
	img.PixelWand.Destroy()
}

func Terminate() {
	imagick.Terminate()
}

func NewImageWand() ImageWand {
	img := ImageWand{
		MagickWand: imagick.NewMagickWand(),
		PixelWand:  imagick.NewPixelWand(),
		DrawWand:   imagick.NewDrawingWand(),
	}
	return img
}

func (img *ImageWand) ResizeImageProcess(data []byte, plan ResizePlan) error {
	err := img.MagickWand.ReadImageBlob(data)
	if err != nil {
		helper.Log.Error("read data failed", err)
		return err
	}
	originWidth := int(img.MagickWand.GetImageWidth())
	originHeight := int(img.MagickWand.GetImageHeight())
	if err = originPictureIsIllegal(originWidth, originHeight); err != nil {
		return err
	}

	o := newResize()
	o.LimitEnlargement = plan.Limit
	o.Background = plan.Color

	if plan.Data != nil {
		if plan.WatermarkProportion != 0 {
			factor, err := factorCalculations(img, plan.Data, float64(plan.WatermarkProportion))
			if err != nil {
				return err
			}
			if err = pictureOverlarge(factor, originWidth, originHeight); err != nil {
				return err
			}
			o.Zoom = factor
			err = img.resize(o)
			if err != nil {
				return err
			}
			return nil
		}
	}

	//proportion zoom
	if plan.Proportion != 0 {
		factor := float64(plan.Proportion) / 100.0
		helper.Log.Info("scaling factor: ", factor)
		o.Zoom = factor
		err = img.resize(o)
		if err != nil {
			return err
		}
		return nil
	}

	o.Width = plan.Width
	o.Height = plan.Height
	switch plan.Mode {
	case "lfit":
		adjustCropTask(&plan, originWidth, originHeight)
		o.Width = plan.Width
		o.Height = plan.Height
		helper.Log.Info("Resize lfit trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "mfit":
		adjustCropTask(&plan, originWidth, originHeight)
		o.Width = plan.Width
		o.Height = plan.Height
		helper.Log.Info("Resize mfit trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "pad":
		o.Pad = true
		helper.Log.Info("Resize pad trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "fixed":
		o.Force = true
		helper.Log.Info("Resize fixed trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "fill":
		o.Crop = true
		helper.Log.Info("Resize fill trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (img *ImageWand) ImageWatermarkProcess(data []byte, plan WatermarkPlan) error {
	err := img.MagickWand.ReadImageBlob(data)
	if err != nil {
		helper.Log.Error("read data failed", err)
		return err
	}
	originWidth := int(img.MagickWand.GetImageWidth())
	originHeight := int(img.MagickWand.GetImageHeight())
	if err = originPictureIsIllegal(originWidth, originHeight); err != nil {
		return err
	}
	w := newWatermark()
	r := newRotate()
	if plan.PictureMask.Image != "" {
		picture := NewImageWand()
		defer picture.Destory()
		err := picture.MagickWand.ReadImageBlob(plan.PictureMask.Data)
		if err != nil {
			helper.Log.Error("open watermark picture file failed")
			return err
		}
		//process rotate
		r.Degrees = plan.RotateDegrees
		r.Background = "transparent"
		helper.Log.Info("Watermark rotate trans params ", r)
		err = picture.rotate(r)
		if err != nil {
			return err
		}

		wmWidth := int(picture.MagickWand.GetImageWidth())
		wmHeight := int(picture.MagickWand.GetImageHeight())

		w.Picture = picture.MagickWand
		w.Transparency = float64(plan.Transparency) / 100.0
		switch plan.Position {
		case NorthWest:
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		case North:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = plan.YMargin
			break
		case NorthEast:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = plan.YMargin
			break
		case West:
			w.XMargin = plan.XMargin
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case Center:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case East:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case SouthWest:
			w.XMargin = plan.XMargin
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		case South:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		case SouthEast:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		default:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = originHeight - plan.YMargin - wmHeight
		}
		helper.Log.Info("Picture watermark trans params ", w, w.Picture)
		err = img.watermark(w)
		if err != nil {
			return err
		}
		return nil
	} else if plan.TextMask.Text != "" {
		t := Text{plan.TextMask.Text,
			plan.TextMask.Color,
			helper.DEFAULT_PIPA_FRONT_PATH + selectTextType(plan.TextMask.Type),
			plan.TextMask.Size,
			plan.TextMask.Shadow,
			plan.RotateDegrees,
			plan.TextMask.Fill,
		}
		w.Picture, err = img.setTextAsPicture(t)
		if err != nil {
			return err
		}
		defer w.Picture.Destroy()

		wmWidth := int(w.Picture.GetImageWidth())
		wmHeight := int(w.Picture.GetImageHeight())

		w.Transparency = float64(plan.Transparency) / 100.0

		switch plan.Position {
		case NorthWest:
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		case North:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = plan.YMargin
			break
		case NorthEast:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = plan.YMargin
			break
		case West:
			w.XMargin = plan.XMargin
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case Center:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case East:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case SouthWest:
			w.XMargin = plan.XMargin
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		case South:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		case SouthEast:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		default:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = originHeight - plan.YMargin - wmHeight
		}
		helper.Log.Info("Text watermark trans params ", w, w.Picture)
		err = img.watermark(w)
		if err != nil {
			return err
		}
		return nil
	} else {
		return ErrInvalidWatermarkProcess
	}
}

func (img *ImageWand) RotateImageProcess(data []byte, plan RotatePlan) error {
	err := img.MagickWand.ReadImageBlob(data)
	if err != nil {
		helper.Log.Error("read data failed", err)
		return err
	}
	originWidth := int(img.MagickWand.GetImageWidth())
	originHeight := int(img.MagickWand.GetImageHeight())

	if originHeight > 4096 || originWidth > 4096 {
		return ErrPictureWidthOrHeightTooLong
	} else if originHeight <= 0 || originWidth <= 0 {
		return ErrPictureWidthOrHeightIsZero
	}

	r := newRotate()
	r.Degrees = plan.Degrees
	r.Background = plan.Color
	helper.Log.Info("Rotate trans params ", r)
	err = img.rotate(r)
	if err != nil {
		return err
	}
	return nil
}

func (img *ImageWand) ReturnData() []byte {
	return img.MagickWand.GetImageBlob()
}
