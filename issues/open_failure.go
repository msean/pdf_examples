package main

// 无法解析合成的pdf，只能通过生成图片的形式粘贴

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

func pdfDownload(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	u := "https://kpserverdev-1251506165.cos.ap-shanghai.myqcloud.com/invoice/upload/20201844_83371_1603161368632.pdf"
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	strem, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ioseeker := io.ReadSeeker(bytes.NewReader(strem))
	pdf := gofpdf.New("P", "pt", "A4", "/MediaBox")
	imp := gofpdi.ImportPageFromStream(pdf, &ioseeker, 1, "/MediaBox")
	pdf.AddPage()
	gofpdi.UseImportedTemplate(pdf, imp, 20, 50, 150, 0)
	if err := pdf.OutputFileAndClose("./out.pdf"); err != nil {
		panic(err)
	}
}
