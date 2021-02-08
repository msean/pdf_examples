package pdf_utils

import (
	"regexp"

	"github.com/jung-kurt/gofpdf"
)

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
		if curWidth > width-2 {
			curWidth = addWidth
			lines = append(lines, string(line))
			line = make([]rune, 0)
		}
		line = append(line, r)
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

// EmojiEscape 正则过滤emoji
func EmojiEscape(text string) string {
	return regexp.MustCompile(`[\x{1F000}-\x{1F9FF}|[\x{2600}-\x{26FF}]`).ReplaceAllString(text, ``)
}

// UnexpectHanEscape 过滤中文状态下的空格
// \x{2006} example "y g n b 的"
func UnexpectHanEscape(text string) string {
	return regexp.MustCompile(`[\x{2006}|\x{202D}]`).ReplaceAllString(text, ` `)
}

func PrintTextEscape(text string) string {
	return EmojiEscape(UnexpectHanEscape(text))
}
