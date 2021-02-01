package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTopMargin(30)

	pdf.AddUTF8Font("SourceHanSansCN", "", "SourceHanSansCN-Bold.ttf")
	pdf.AddUTF8Font("SourceHanSansCN", "B", "SourceHanSansCN-Normal.ttf")
	pdf.SetFont("SourceHanSansCN", "", 12)

	//设置页眉
	pdf.SetHeaderFuncMode(func() {
		pdf.SetY(10)
		pdf.CellFormat(0, 9, "examples", "0", 0, "C", false, 0, "")
	}, true)

	//设置页脚
	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.CellFormat(0, 10, fmt.Sprintf("当前第 %d 页", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	pdf.AddPage()
	pdf.AddPage()
	if err := pdf.OutputFileAndClose("../assets/add_footer.pdf"); err != nil {
		panic(err.Error())
	}
}
