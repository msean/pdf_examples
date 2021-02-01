package main

import "github.com/jung-kurt/gofpdf"

func main() {
	//设置页面参数
	pdf := gofpdf.New("P", "mm", "A4", "")

	//添加一页
	pdf.AddPage()

	//将图片放入到 pdf 文档中
	//ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	pdf.ImageOptions("../assets/picture.jpg", 20, 20, 120, 120, false, gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: false}, 0, "")

	pdf.AddPage()
	pdf.TransformBegin()
	pdf.TransformRotate(90, 80, 80)

	pdf.ImageOptions(
		"../assets/picture.jpg",
		20, 20,
		120, 120,
		false,
		gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: false},
		0,
		"",
	)
	pdf.TransformEnd()

	if err := pdf.OutputFileAndClose("../assets/rotate.pdf"); err != nil {
		panic(err.Error())
	}

}
