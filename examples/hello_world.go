package main

import "github.com/jung-kurt/gofpdf"

func main() {
	// 创建pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	// 增加一页
	pdf.AddPage()
	// 添加字体
	pdf.AddUTF8Font("SourceHanSansCN", "", "SourceHanSansCN-Bold.ttf")
	pdf.AddUTF8Font("SourceHanSansCN", "B", "SourceHanSansCN-Normal.ttf")
	// 设置字体
	pdf.SetFont("SourceHanSansCN", "B", 8)
	// 绘制单元格
	// CellFormat(w, h float64, txtStr, borderStr string, ln int,alignStr string, fill bool, link int, linkStr string)
	pdf.CellFormat(100, 100, "Hello World", "1", 0, "CM", false, 0, "")
	if err := pdf.OutputFileAndClose("../assets/create_pdf.pdf"); err != nil {
		panic(err)
	}
}
