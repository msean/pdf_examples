#### 一、pdf第三方库
> github.com/jung-kurt/gofpdf

##### 缺点：
* 对中文支持度不是太友好
* 无法支持表情字符,遇到有表情符号的字符将直接panic导致程序崩溃


#### 二、pdf基本操作
##### 1、hello world
``` go
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
	if err := pdf.OutputFileAndClose("pdfs/create_pdf.pdf"); err != nil {
		panic(err)
	}
}
```
* CellFormat 绘制单元格几个重要函数参数

|  参数名   | 类型  | 含义  | 备注  |
|  ----    | ---- | ----  | ----  |
| w  | float64 | 内容所占宽度 |  |
|  h  | float64 | 内容所占高度 |  |
|  txtStr  | float64 | 内容 |   |
|  borderStr  | float64 | 是否有边框 | 1 全边框 /L/R/""无边框|
|  alignStr  | float64 | 内容对齐方式 | 默认CM 垂直居中, L 左对齐|
|  fill  | bool | 单元格 | 是否填充背景颜色|

* 目前中文字体格式只支持ttf格式(格式有ttf, otf, ttc)
> 几种字体之间的区别 https://jingyan.baidu.com/article/5d6edee2fe14f299eadeec1c.html
* 思源true-type 字体大全
> https://github.com/Pal3love/Source-Han-TrueType


##### 2、添加图片
``` go
func main() {
	//设置页面参数
	pdf := gofpdf.New("P", "mm", "A4", "")

	//添加一页
	pdf.AddPage()

	//将图片放入到 pdf 文档中
	//ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	pdf.ImageOptions("pdfs/picture.jpg", 0, 0, 0, 0, false, gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: false}, 0, "")

	if err := pdf.OutputFileAndClose("pdfs/add_picture.pdf"); err != nil {
		panic(err.Error())
	}
}
```

#### 3、设置页眉页脚
```go
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
	if err := pdf.OutputFileAndClose("pdfs/add_footer.pdf"); err != nil {
		panic(err.Error())
	}
}
```
#### 4、pdf粘贴
```go
package main

import (
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

func main() {
	srcFilePath := "pdfs/add_footer.pdf"
	dstFilePath := "pdfs/paster.pdf"
	pdf := gofpdf.New("P", "pt", "A4", "")
	imp := gofpdi.NewImporter()
	tpl := imp.ImportPage(pdf, srcFilePath, 1, "/MediaBox")
	pageSizes := imp.GetPageSizes()
	nrPages := len(imp.GetPageSizes())
	// add all pages from template pdf
	for i := 1; i <= nrPages; i++ {
		pdf.AddPage()
		if i > 1 {
			tpl = imp.ImportPage(pdf, srcFilePath, i, "/MediaBox")
		}
		imp.UseImportedTemplate(pdf, tpl, 0, 0, pageSizes[i]["/MediaBox"]["w"], pageSizes[i]["/MediaBox"]["h"])
	}

	if err := pdf.OutputFileAndClose(dstFilePath); err != nil {
		panic(err)
	}
}
```
#### 5、旋转
```go
func main() {
	//设置页面参数
	pdf := gofpdf.New("P", "mm", "A4", "")

	//添加一页
	pdf.AddPage()

	//将图片放入到 pdf 文档中
	//ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	pdf.ImageOptions("pdfs/picture.jpg", 20, 20, 120, 120, false, gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: false}, 0, "")

	pdf.AddPage()
	pdf.TransformBegin()
	// 以坐标(x, y)(80, 80)为中心点旋转90度
	pdf.TransformRotate(90, 80, 80)

	pdf.ImageOptions(
		"pdfs/picture.jpg",
		20, 20,
		120, 120,
		false,
		gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: false},
		0,
		"",
	)
	pdf.TransformEnd()

	if err := pdf.OutputFileAndClose("pdfs/rotate.pdf"); err != nil {
		panic(err.Error())
	}
}
```

#### 三、报销单打印遇到的问题
##### 1、特许字符和表情包字符支持度不友好
* 表情包字符会直接引起程序奔溃，过滤表情字符
```go
regexp.MustCompile(`[\x{1F000}-\x{1F9FF}|[\x{2600}-\x{26FF}]`).ReplaceAllString(text, ``)
```
* 特殊字符打印会表现形式为 ·口· 字（go test ./test -v -run Test_CreatePdf）
```go
regexp.MustCompile(`[\x{2006}|\x{202D}]`).ReplaceAllString(text, ` `)
```

##### 2、库中自带的分行函数(SplitLines)对中文每行的字不一样,碰到需要换行的需要自行判断
```go
// DynamicCol 根据Grid的内容计算所需占用的列数
func DynamicCol(width float64, text string, pdf *gofpdf.Fpdf) (rows int, lines []string) {
	curWidth := 0.0
	lines = make([]string, 0)
	line := make([]rune, 0)
	runes := []rune(text)
	linebreak := []rune("\n")[0]
	for i, r := range runes {
		// 换行符
		if r == linebreak {
			lines = append(lines, string(line))
			curWidth = 0.0
			line = make([]rune, 0)
			continue
		}
		addWidth := 0.0
		addWidth += pdf.GetStringWidth(string(r))
		curWidth += addWidth
		// 超出宽度，自动换行
		if curWidth > width{
			curWidth = addWidth
			lines = append(lines, string(line))
			line = make([]rune, 0)
		}
		line = append(line, r)
		// 最后一个字符
		if i == len(runes)-1 && len(line) > 0 {
			lines = append(lines, string(line))
		}
	}
	rows = len(lines)
	if rows == 0 {
		lines = []string{""}
		rows = 1
	}
	return
}
```
##### 3、绘制一行多列的单元格，某些列可能会换行，所有的列都需要垂直居中
* 找出行数最高的列，根据行数最高的列来调整所有列的宽度，实现方见Col类

##### 4、绘制多行多列的单元格，绘制单元格过程中可能出现翻页的情况，需要将标题栏带入下一页
* 实现方法见ColWithDiving类
* 局限性：
1、只能适应一列中有1或者n行(n是固定的)的情况 
2、单行的高度不能大于多行的总和

##### 5、根据图片的大小尺寸和PDF的大小尺寸适应旋转
```go
func rotateImg(m image.Image) (dstm image.Image) {
	img := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			img.Set(x, m.Bounds().Max.X-y, m.At(y, x))
		}
	}
	dstm = img
	return
}
```
##### 6、头部旋转
* 寻找合适的旋转中心点(利用正方形以中心旋转90/180/270所占空间位置不变的原理)
##### 7、每页显示总页数
* 在画完之后重新创建一张空白的pdf将原来画好的pdf粘贴到新的pdf上(此时原pdf的总页数可以获取)，并设置新pdf的每一页的页眉

#### 四、related
1、pdf programma
> https://planetpdf.com/planetpdf/pdfs/pdf2k/01W/rosenthol_intro2pdfprog.pdf

2、emoji表情符字符编码对照表
> https://apps.timwhitlock.info/emoji/tables/unicode#block-6c-other-additional-symbols

3、examples code
> https://github.com/msean/pdf_examples
