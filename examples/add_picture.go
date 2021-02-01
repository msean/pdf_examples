package main

import (
	"github.com/jung-kurt/gofpdf"
)

func main() {
	//设置页面参数
	pdf := gofpdf.New("P", "mm", "A4", "")

	//添加一页
	pdf.AddPage()

	//将图片放入到 pdf 文档中
	//ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	pdf.ImageOptions("../assets/picture.jpg", 0, 0, 0, 0, false, gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: false}, 0, "")

	if err := pdf.OutputFileAndClose("../assets/add_picture.pdf"); err != nil {
		panic(err.Error())
	}
}
