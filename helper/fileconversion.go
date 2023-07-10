package helper

import (
	"errors"
	"fmt"
	"github.com/raykov/gofpdf"
	"github.com/raykov/mdtopdf/document"
	"github.com/raykov/mdtopdf/renderer"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/text/encoding/charmap"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func MD2PDF(command *Command) error {
	if command.Args == nil || len(command.Args) == 0 {
		return errors.New("no input file")
	}
	inputFile := command.Args[0]
	tempFile := inputFile
	// Remove any extension
	if strings.Contains(tempFile, ".") {
		tempFile = tempFile[:strings.LastIndex(tempFile, ".")]
	}
	outputFile := tempFile + ".pdf"
	if len(command.Args) > 1 {
		outputFile = command.Args[1]
	}
	stack := GetStack()
	if inputFile[0] != '/' {
		inputFile = filepath.Join(stack.ActivePath, inputFile)
	}
	if outputFile[0] != '/' {
		outputFile = filepath.Join(stack.ActivePath, outputFile)
	}
	md, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer md.Close()

	pdf, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer pdf.Close()

	pageNumExtension := func(pdf *gofpdf.Fpdf) {
		pdf.SetFooterFunc(func() {
			left, _, right, bottom := pdf.GetMargins()
			width, height := pdf.GetPageSize()
			fontSize := 12.0

			pNum := fmt.Sprint(pdf.PageNo())
			pdf.SetXY(width-left/2-pdf.GetStringWidth(pNum), height-bottom/2)
			pdf.SetFontSize(fontSize)
			pdf.SetTextColor(200, 200, 200)
			pdf.SetFontStyle("B")
			pdf.SetRightMargin(0)
			pdf.Write(fontSize, pNum)
			pdf.SetRightMargin(right)
		})
	}

	err = convertUTF8(md, pdf, pageNumExtension)
	if err != nil {
		return err
	}

	return nil
}

// Convert your Markdown to PDF
func convertUTF8(r io.Reader, w io.Writer, extensions ...func(*gofpdf.Fpdf)) error {
	md, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	out := make([]byte, 0)

	for _, r := range string(md) {
		if e, ok := charmap.ISO8859_1.EncodeRune(r); ok {
			out = append(out, e)
		}
	}

	markdown := goldmark.New(
		goldmark.WithRenderer(renderer.NewRenderer()),
		goldmark.WithExtensions(
			extension.NewTable(),
			extension.Strikethrough,
		),
	)

	pdf := gofpdf.New("P", "pt", "A4", ".")

	for _, extension := range extensions {
		extension(pdf)
	}

	pdf.AddPage()

	d := document.NewDocument(pdf, document.DefaultStyle)

	if err = markdown.Convert(out, d); err != nil {
		return err
	}

	err = pdf.Output(w)
	if err != nil {
		return err
	}

	return nil
}
