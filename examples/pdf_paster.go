package main

import (
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

func main() {
	srcFilePath := "../assets/add_footer.pdf"
	dstFilePath := "../assets/paster.pdf"
	pdf := gofpdf.New("P", "pt", "A4", "")
	imp := gofpdi.NewImporter()
	tpl := imp.ImportPage(pdf, srcFilePath, 1, "/MediaBox")
	pageSizes := imp.GetPageSizes()
	nrPages := len(imp.GetPageSizes())
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
