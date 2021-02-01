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
	// 绘制单元格
	pdf.SetXY(5, 20)
	pdf_utils.NewCol(pdf, pdf_utils.DefaultLeftX, pdf_utils.DefaultLineHeight, []pdf_utils.Grid{
		{Text: "行程", Width: 10, Fill: true},   //行程
		{Text: "初始地", Width: 25, Fill: true},  // 初始地
		{Text: "目的地", Width: 25, Fill: true},  //目的地
		{Text: "出行方式", Width: 15, Fill: true}, // 出行方式
		{Text: "开始时间", Width: 30, Fill: true}, // 开始时间
		{Text: "结束时间", Width: 30, Fill: true}, // 结束时间
		{Text: "出差时间", Width: 25, Fill: true}, // 出差时间
		{Text: "相关客户", Width: 40, Fill: true}, // 相关客户
	}...).Draw()
	pdf_utils.NewCol(pdf, pdf_utils.DefaultLeftX, pdf_utils.DefaultLineHeight, []pdf_utils.Grid{
		{Text: "北京", Width: 10},   //行程
		{Text: "北京通州", Width: 25}, // 初始地
		{Text: "北京朝阳", Width: 25}, //目的地
		{Text: "飞机", Width: 15},   // 出行方式
		{Text: "开始时间", Width: 30}, // 开始时间
		{Text: "这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材", Width: 30, IsDynamic: true}, // 结束时间
		{Text: "出差时间", Width: 25}, // 出差时间
		// {Text: "这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材", Width: 40, IsDynamic: true, Bold: true}, // 相关客户
		{Text: "这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材这个传闻会是闲谈专栏作家的又一素材", Width: 40, IsDynamic: true},
	}...).Draw()
	if err := pdf.OutputFileAndClose("../assets/col.pdf"); err != nil {
		panic(err)
	}
}
