package pdf_utils

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type (
	Grid struct {
		IsDynamic    bool //根据内容可能需要动态高度
		Text         string
		Width        float64
		Alian        string
		NewlineAlian string // 换行align
		Bold, Fill   bool   // 是否需要加粗 / 是否需要设置背景颜色
	}
	Col struct {
		Grids                         []Grid
		pdf                           *gofpdf.Fpdf
		startx, BaseHeight, fixHeight float64
	}
	ComplexGrid struct {
		Ismuti    bool
		lineTexts []string
		width     float64
		alian     string
		bold      bool
		fill      bool
	}
	ColWithDiving struct {
		multiGrids []*ComplexGrid
		pdf        *gofpdf.Fpdf
		lines      int
		baseHeight float64
		startx     float64
		flipFunc   func() // 翻页执行函数
	}
)

// NewCol NewCol
// 默认居中对齐/换行也是居中对齐
func NewCol(pdf *gofpdf.Fpdf, startx float64, baseHeight float64, grids ...Grid) *Col {
	col := new(Col)
	col.pdf = pdf
	col.Grids = make([]Grid, 0)
	col.startx = startx
	if baseHeight == 0 {
		baseHeight = DefaultLineHeight
	}
	col.BaseHeight = baseHeight
	col.fixHeight = baseHeight
	for _, grid := range grids {
		if grid.Alian == "" {
			// 默认居中
			grid.Alian = "CM"
		}
		if grid.NewlineAlian == "" {
			grid.NewlineAlian = "CM"
		}
		col.Grids = append(col.Grids, grid)
	}
	return col
}

func (col *Col) MaxGridLines() (g *Grid) {
	var firstDynamicMark bool
	for _, grid := range col.Grids {
		curgrid := grid
		if grid.IsDynamic {
			if !firstDynamicMark {
				g = &curgrid
				firstDynamicMark = true
			} else {
				if float64(len(grid.Text))/(grid.Width) > float64(len(g.Text))/(g.Width) {
					g = &curgrid
				}
			}
		}
	}
	return
}

// NeedFix 整行是否需要做动态调整
func (col *Col) NeedFix() (bool, int) {
	g := col.MaxGridLines()
	if g != nil {
		// 需要做调整
		lines, _ := DynamicCol(g.Width, g.Text, col.pdf)
		if lines > 1 {
			col.fixHeight = 0.0
			for i := 0; i < lines; i++ {
				col.fixHeight += lineHeightSpace
			}
			col.fixHeight += mutilLineTopMargin + mutilLineBottomMargin
			return true, lines
		}
	}
	return false, 1
}

func setGridStartY(colHeight, startY, textHeight float64, lines int) float64 {
	textTakeHeight := float64(lines) * textHeight
	spaceHeight := colHeight - textTakeHeight - mutilLineTopMargin - mutilLineBottomMargin
	return startY + spaceHeight/2 + mutilLineTopMargin
}

// Draw 绘画整条col
func (col *Col) Draw() {
	fix, _ := col.NeedFix()
	fmt.Println(">>>>>>>>>>", col.fixHeight)
	gridxStart := col.startx
	starty := col.pdf.GetY()
	for _, grid := range col.Grids {
		// 加粗
		if grid.Bold {
			SetTextColor(col.pdf, BordRGB)
			SetFont(col.pdf, DefaultBoldFont, DefaultFontSize)
		}
		col.pdf.SetXY(gridxStart, starty)
		gridxStart += grid.Width
		if fix && grid.IsDynamic {
			lines, linesText := DynamicCol(grid.Width, grid.Text, col.pdf)
			if lines > 1 {
				tmpX, tmpY := col.pdf.GetX(), col.pdf.GetY()
				// 设置起始高度 是垂直居中
				gridStartY := setGridStartY(col.RealHeight(), tmpY, lineHeightSpace, lines)
				col.pdf.CellFormat(grid.Width, col.fixHeight, "", "1", 0, grid.Alian, grid.Fill, 0, "")
				col.pdf.SetXY(tmpX, gridStartY)
				for index, lineText := range linesText {
					col.pdf.SetX(tmpX)
					alian := grid.Alian
					if index != 0 {
						alian = grid.NewlineAlian
					}
					col.pdf.CellFormat(grid.Width, lineHeightSpace, lineText, "", 0, alian, grid.Fill, 0, "")
					if index != lines-1 {
						col.pdf.Ln(lineHeightSpace)
					} else {
						col.pdf.Ln(mutilLineBottomMargin)
					}
				}
			} else {
				col.pdf.CellFormat(grid.Width, col.fixHeight, grid.Text, "1", 0, grid.Alian, grid.Fill, 0, "")
			}
		} else {
			col.pdf.CellFormat(grid.Width, col.fixHeight, grid.Text, "1", 0, grid.Alian, grid.Fill, 0, "")
		}
		// 还原默认设置
		if grid.Bold {
			SetTextColor(col.pdf, DefaultRGB)
			SetFont(col.pdf, DefaultFont, DefaultFontSize)
		}
	}
	col.pdf.SetY(starty + col.RealHeight())
}

// RealHeight 获取实际的高度
func (col *Col) RealHeight() float64 {
	if col.fixHeight > col.BaseHeight {
		return col.fixHeight
	}
	return col.BaseHeight
}

// NewComplexGrid NewComplexGrid
// alian 内容对齐方式
// bold 内容是否需要加粗
func NewComplexGrid(texts []string, width float64, alian string, bold bool, fill bool) *ComplexGrid {
	var ismuti bool
	if len(texts) > 1 {
		ismuti = true
	}
	return &ComplexGrid{
		Ismuti:    ismuti,
		width:     width,
		lineTexts: texts,
		alian:     alian,
		bold:      bold,
		fill:      fill,
	}
}

// NewDivingCol NewDivingCol
// 设计缺陷：保证单行格不能有换行
func NewDivingCol(pdf *gofpdf.Fpdf, startx, baseHeight float64, flipFunc func(), complexGrids ...*ComplexGrid) (col *ColWithDiving, err error) {
	col = new(ColWithDiving)
	col.lines = 1
	for _, complexGrid := range complexGrids {
		lines := len(complexGrid.lineTexts)
		if lines > col.lines {
			col.lines = lines
		}
		if lines != 1 && lines != col.lines {
			err = fmt.Errorf("NewDivingColErr")
			return
		}
	}
	// 假如是行数为1
	if col.lines == 1 {
		for _, complexGrid := range complexGrids {
			complexGrid.Ismuti = true
		}
	}
	col.pdf = pdf
	col.multiGrids = make([]*ComplexGrid, 0)
	col.startx = startx
	col.baseHeight = baseHeight
	col.multiGrids = complexGrids
	col.flipFunc = flipFunc
	return
}

// SetGridX Draw
func (c *ColWithDiving) SetGridX(seq int) float64 {
	startx := c.startx
	for index, grid := range c.multiGrids {
		if index == seq {
			return startx
		}
		startx += grid.width
	}
	return startx
}

// Draw Draw
func (c *ColWithDiving) Draw() {
	curLine, curY := 1, c.pdf.GetY()
	cliFuncHeight := 0.0
	sigleGridY, singleGridH := c.pdf.GetY(), 0.0
	for curLine <= c.lines {
		lindex := curLine - 1
		addY := c.baseHeight
		// 计算该行所占用的高度
		for _, grid := range c.multiGrids {
			if grid.Ismuti && len(grid.lineTexts) >= lindex {
				rows, _ := DynamicCol(grid.width, grid.lineTexts[lindex], c.pdf)
				h := lineHeightSpace * float64(rows)
				if rows > 1 {
					h += mutilLineTopMargin + mutilLineBottomMargin
				}
				if h > addY {
					addY = h
				}
			}
		}
		// 超过底部 则需要换页 将非multi的格画好
		if (curY + addY) > DefaultBottomY {
			curY = DefaultTopY
			if singleGridH != 0.0 {
				for i, mutiGrid := range c.multiGrids {
					if !mutiGrid.Ismuti {
						// 加粗
						if mutiGrid.bold {
							SetTextColor(c.pdf, BordRGB)
							SetFont(c.pdf, DefaultBoldFont, DefaultFontSize)
						}
						gridStartx := c.SetGridX(i)
						c.pdf.SetXY(gridStartx, sigleGridY)
						var text string
						if len(mutiGrid.lineTexts) > 0 {
							text = mutiGrid.lineTexts[0]
						}
						rows, textLines := DynamicCol(mutiGrid.width, text, c.pdf)
						if rows > 1 {
							c.pdf.SetX(gridStartx)
							c.pdf.CellFormat(mutiGrid.width, singleGridH, "", "1", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
							// 设置起始y坐标
							c.pdf.SetXY(gridStartx, setGridStartY(singleGridH, c.pdf.GetY(), lineHeightSpace, rows))
							for index, lineText := range textLines {
								if index == 0 {
									c.pdf.Ln(mutilLineTopMargin)
								}
								c.pdf.SetX(gridStartx)
								c.pdf.CellFormat(mutiGrid.width, lineHeightSpace, lineText, "0", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
								if index != rows-1 {
									c.pdf.Ln(lineHeightSpace)
								} else {
									c.pdf.Ln(mutilLineBottomMargin)
								}
							}
						} else {
							c.pdf.SetX(gridStartx)
							var text string
							if len(textLines) > 0 {
								text = textLines[0]
							}
							c.pdf.CellFormat(mutiGrid.width, singleGridH, text, "1", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
						}
						// 加粗后还原
						if mutiGrid.bold {
							SetTextColor(c.pdf, BordRGB)
							SetFont(c.pdf, DefaultFont, DefaultFontSize)
						}
					}
				}
			}
			// 执行翻页函数
			c.pdf.AddPage()
			c.pdf.SetY(DefaultTopY)
			if c.flipFunc != nil {
				cys := c.pdf.GetY()
				c.flipFunc()
				cye := c.pdf.GetY()
				cliFuncHeight += cye - cys
				curY += cliFuncHeight
			}
			singleGridH = 0.0
			sigleGridY = c.pdf.GetY()
		}

		// 画multi格子
		for i, mutiGrid := range c.multiGrids {
			if mutiGrid.bold {
				SetTextColor(c.pdf, BordRGB)
				SetFont(c.pdf, DefaultBoldFont, DefaultFontSize)
			}
			if mutiGrid.Ismuti {
				rows, textLines := DynamicCol(mutiGrid.width, mutiGrid.lineTexts[lindex], c.pdf)
				gridStartx := c.SetGridX(i)
				if rows > 1 {
					c.pdf.SetX(gridStartx)
					c.pdf.CellFormat(mutiGrid.width, addY, "", "1", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
					c.pdf.SetXY(gridStartx, setGridStartY(addY, c.pdf.GetY(), lineHeightSpace, rows))
					for index, lineText := range textLines {
						c.pdf.SetX(gridStartx)
						// 设置起始y坐标
						c.pdf.CellFormat(mutiGrid.width, lineHeightSpace, lineText, "0", 0, "L", mutiGrid.fill, 0, "")
						if index != rows-1 {
							c.pdf.Ln(lineHeightSpace)
						} else {
							c.pdf.Ln(mutilLineBottomMargin)
						}
					}
				} else {
					c.pdf.SetXY(gridStartx, c.pdf.GetY())
					var text string
					if len(textLines) > 0 {
						text = textLines[0]
					}
					c.pdf.CellFormat(mutiGrid.width, addY, text, "1", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
				}
				if mutiGrid.bold {
					SetTextColor(c.pdf, DefaultRGB)
					SetFont(c.pdf, DefaultFont, DefaultFontSize)
				}
			}
			c.pdf.SetY(curY)
		}
		singleGridH += addY
		curY += addY
		cliFuncHeight = 0.0
		if curLine == c.lines {
			for i, mutiGrid := range c.multiGrids {
				if mutiGrid.bold {
					SetTextColor(c.pdf, BordRGB)
					SetFont(c.pdf, DefaultBoldFont, DefaultFontSize)
				}
				if !mutiGrid.Ismuti {
					gridStartx := c.SetGridX(i)
					c.pdf.SetXY(gridStartx, sigleGridY)
					var text string
					if len(mutiGrid.lineTexts) > 0 {
						text = mutiGrid.lineTexts[0]
					}
					rows, textLines := DynamicCol(mutiGrid.width, text, c.pdf)
					if rows > 1 {
						c.pdf.SetX(gridStartx)
						c.pdf.CellFormat(mutiGrid.width, singleGridH, "", "1", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
						// 设置起始y坐标
						c.pdf.SetXY(gridStartx, setGridStartY(singleGridH, c.pdf.GetY(), lineHeightSpace, rows))
						for index, lineText := range textLines {
							c.pdf.SetX(gridStartx)
							c.pdf.CellFormat(mutiGrid.width, lineHeightSpace, lineText, "0", 0, "L", mutiGrid.fill, 0, "")
							if index != rows-1 {
								c.pdf.Ln(lineHeightSpace)
							} else {
								c.pdf.Ln(mutilLineBottomMargin)
							}
						}
					} else {
						c.pdf.SetX(gridStartx)
						var text string
						if len(textLines) > 0 {
							text = textLines[0]
						}
						c.pdf.CellFormat(mutiGrid.width, singleGridH, text, "1", 0, mutiGrid.alian, mutiGrid.fill, 0, "")
					}
				}
				if mutiGrid.bold {
					SetTextColor(c.pdf, DefaultRGB)
					SetFont(c.pdf, DefaultFont, DefaultFontSize)
				}
			}
		}
		curLine++
	}
	c.pdf.SetXY(c.startx, curY)
}
