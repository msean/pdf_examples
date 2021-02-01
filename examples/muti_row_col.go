package main

import (
	"pdf/pdf_utils"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// 创建pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("SourceHanSansCN", "B", "SourceHanSansCN-Bold.ttf")
	pdf.AddUTF8Font("SourceHanSansCN", "", "SourceHanSansCN-Normal.ttf")
	pdf.SetFont("SourceHanSansCN", "", pdf_utils.DefaultFontSize)
	pdf_utils.SetBackgroudColor(pdf, pdf_utils.DefaultBackgroudRGB)
	pdf.SetXY(5, pdf_utils.DefaultTopY)
	BarFunc := func() {
		pdf_utils.NewCol(pdf, pdf_utils.DefaultLeftX, pdf_utils.DefaultLineHeight, []pdf_utils.Grid{
			{Text: "行程", Width: 10, Fill: true},      //行程
			{Text: "初始地-目的地", Width: 25, Fill: true}, // 初始地
			{Text: "出行方式", Width: 15, Fill: true},    // 出行方式
			{Text: "开始时间", Width: 30, Fill: true},    // 开始时间
			{Text: "结束时间", Width: 30, Fill: true},    // 结束时间
			{Text: "出差时间", Width: 50, Fill: true},    // 出差时间
			{Text: "相关客户/备注", Width: 40, Fill: true}, // 相关客户
		}...).Draw()
	}
	BarFunc()
	for i := 0; i < 100; i++ {
		complexGrids := []*pdf_utils.ComplexGrid{
			pdf_utils.NewComplexGrid([]string{"北京"}, 10, "CM", false, false),
			pdf_utils.NewComplexGrid([]string{"北京通州", "北京朝阳"}, 25, "CM", false, false),
			pdf_utils.NewComplexGrid([]string{"飞机"}, 15, "CM", false, false),
			pdf_utils.NewComplexGrid([]string{"开始时间"}, 30, "CM", false, false),
			pdf_utils.NewComplexGrid([]string{"结束时间"}, 30, "CM", false, false),
			pdf_utils.NewComplexGrid([]string{"出差时间"}, 50, "CM", false, false),
			pdf_utils.NewComplexGrid([]string{"这个传闻会是闲谈专栏作家的又一素材这个这个传闻会是闲谈专栏作家的又一素材这个", "这个传闻会是闲谈专栏作家的又一素材这个这个传闻会是闲谈专栏作家的又一素材这个"}, 40, "CM", false, false),
		}
		if col, e := pdf_utils.NewDivingCol(pdf, pdf_utils.DefaultLeftX, pdf_utils.DefaultLineHeight, BarFunc, complexGrids...); e == nil {
			col.Draw()
		}
	}
	if err := pdf.OutputFileAndClose("../assets/muti_row_col.pdf"); err != nil {
		panic(err)
	}
}
