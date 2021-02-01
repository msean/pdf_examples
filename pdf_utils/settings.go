package pdf_utils

import "github.com/jung-kurt/gofpdf"

type (
	RGB struct {
		R, G, B int
	}
	Font struct {
		Family, Style, File string
	}
)

// GetPageSize
const (
	PageWidth         = 210.0
	PageHeight        = 297.0
	PageContentWidth  = 200.0
	PageContentHeight = PageHeight - DefaultTopY - 15.0
)

const (
	DefaultLineHeight     = 6.5
	DefaultTopY           = 42.0
	DefaultLeftX          = 5.0
	DefaultBottomY        = 275.0
	DefaultFontSize       = 9
	DefaultHeaderFontSize = 9.5
	BlockSpace            = 4.0
)

const (
	lineHeightSpace       = DefaultLineHeight - 1.6 // 换行两行之间的字体间距
	mutilLineTopMargin    = 1.4
	mutilLineBottomMargin = 1.4
)

var (
	DefaultRGB = RGB{
		R: 90,
		G: 90,
		B: 90,
	}
	BordRGB = RGB{
		R: 0,
		G: 0,
		B: 0,
	}
	DefaultBackgroudRGB = RGB{
		R: 240,
		G: 240,
		B: 240,
	}
)

var (
	DefaultFont = Font{
		Family: "SourceHanSansCN",
		Style:  "",
		File:   "SourceHanSansCN-Normal.ttf",
	}
	DefaultBoldFont = Font{
		Family: "SourceHanSansCN",
		Style:  "B",
		File:   "SourceHanSansCN-Bold.ttf",
	}
)

// SetTextColor 设置字体颜色
func SetTextColor(pdf *gofpdf.Fpdf, rgb RGB) {
	pdf.SetTextColor(rgb.R, rgb.G, rgb.B)
}

// SetBackgroudColor 设置背景颜色
func SetBackgroudColor(pdf *gofpdf.Fpdf, rgb RGB) {
	pdf.SetFillColor(rgb.R, rgb.G, rgb.B)
}

// AddUTF8Font 设置字体颜色
func AddUTF8Font(pdf *gofpdf.Fpdf, font Font) {
	pdf.AddUTF8Font(font.Family, font.Style, font.File)
}

// SetFont 设置字体颜色
func SetFont(pdf *gofpdf.Fpdf, font Font, size float64) {
	pdf.SetFont(font.Family, font.Style, size)
}
