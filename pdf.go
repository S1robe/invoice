package main

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	subtotalLabel   = "Subtotal"
	discountLabel   = "Discount"
	taxLabel        = "Tax"
	totalLabel      = "Total"
	totalHoursLabel = "Total Hours"
)

func handleRemainingSpace(pdf *gopdf.GoPdf, consumed float64) {
  if(remainingSpace - consumed < 0){
    writeFooter(pdf, file.Id)
    pdf.AddPage()
    remainingSpace = maxHeight
  }
  remainingSpace -= consumed
}

func writeLogo(pdf *gopdf.GoPdf, logo string, from string) {
  if logo != "" {
		width, height := getImageDimension(logo)
		scaledWidth := 100.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
    
    handleRemainingSpace(pdf, float64(scaledHeight))
    _ = pdf.Image(logo, pdf.GetX(), pdf.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight})
	  handleRemainingSpace(pdf, scaledHeight + 24)
    pdf.Br(scaledHeight + 24)
	}
	pdf.SetTextColor(55, 55, 55)

	formattedFrom := strings.ReplaceAll(from, `\n`, "\n")
	fromLines := strings.Split(formattedFrom, "\n")
  // # of lines * font size - 1 cause first is 12, minus the line breaks at the end
  handleRemainingSpace(pdf, float64(len(fromLines) * 25 - 25 + 30 + 21))

	for i := 0; i < len(fromLines); i++ {
		if i == 0 {
			_ = pdf.SetFont("Inter", "", 12)
			_ = pdf.Cell(nil, fromLines[i])
			pdf.Br(18)
		} else {
			_ = pdf.SetFont("Inter", "", 10)
			_ = pdf.Cell(nil, fromLines[i])
			pdf.Br(15)
		}
	}
	pdf.Br(21)

  handleRemainingSpace(pdf, 36)

	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX(), pdf.GetY(), 260, pdf.GetY())
	pdf.Br(36)
}

func writeTitle(pdf *gopdf.GoPdf, title, id, date string) {
	handleRemainingSpace(pdf, 24 + 36 + 12 + 48)
  _ = pdf.SetFont("Inter-Bold", "", 24)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.Cell(nil, title)
	pdf.Br(36)
	_ = pdf.SetFont("Inter", "", 12)
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, "#")
	_ = pdf.Cell(nil, id)
	pdf.SetTextColor(150, 150, 150)
	_ = pdf.Cell(nil, "  ·  ")
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, date)
	pdf.Br(48)
}

func writeDueDate(pdf *gopdf.GoPdf, due string) {
  handleRemainingSpace(pdf, 9 + 11 + 12)
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "Due Date")
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(11)
	pdf.SetX(amountColumnOffset - 15)
	_ = pdf.Cell(nil, due)
	pdf.Br(12)
}

func writeBillTo(pdf *gopdf.GoPdf, to string) {
  formattedTo := strings.ReplaceAll(to, `\n`, "\n")
  toLines := strings.Split(formattedTo, "\n")
  
  handleRemainingSpace(pdf, float64( 9 + 18 + (len(toLines) * 25 - 25 + 35 + 64)))


  pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 9)
	_ = pdf.Cell(nil, "BILL TO")
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)


	for i := 0; i < len(toLines); i++ {
		if i == 0 {
			_ = pdf.SetFont("Inter", "", 15)
			_ = pdf.Cell(nil, toLines[i])
			pdf.Br(20)
		} else {
			_ = pdf.SetFont("Inter", "", 10)
			_ = pdf.Cell(nil, toLines[i])
			pdf.Br(15)
		}
	}
	pdf.Br(64)
}

//Itemization section
func writeHeaderRow(pdf *gopdf.GoPdf) {
  handleRemainingSpace(pdf, 9 + 24)

	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "ITEM")
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, "QTY")
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "RATE")
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, "AMOUNT")
	pdf.Br(24)
}

func writeNotes(pdf *gopdf.GoPdf, notes string) {
	handleRemainingSpace(pdf, (maxHeight - 600))
  pdf.SetY(600)
  formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
  notesLines := strings.Split(formattedNotes, "\n")

  // May not need this...
  //handleRemainingSpace(pdf, float64(9 + 18 + len(notesLines) * (9+15) + 48))
  
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "NOTES")
	pdf.Br(18)
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(0, 0, 0)


	for i := 0; i < len(notesLines); i++ {
		_ = pdf.Cell(nil, notesLines[i])
		pdf.Br(15)
	}

	pdf.Br(48)
}

func writeFooter(pdf *gopdf.GoPdf, id string) {
  // Nothing is ever written below 800 per the handleRemainingSpace function
  pdf.SetY(800)

	_ = pdf.SetFont("Inter", "", 10)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, id)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
	pdf.Br(48)
}

func writeRow(pdf *gopdf.GoPdf, item string, quantity float64, rate float64) {
  handleRemainingSpace(pdf, 11+24)
	_ = pdf.SetFont("Inter", "", 11)
	pdf.SetTextColor(0, 0, 0)

	// Round to 2 dcimal places
	total := float64(quantity) * rate
	amount := strconv.FormatFloat(total, 'f', 2, 64)

	_ = pdf.Cell(nil, item)
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, strconv.FormatFloat(quantity, 'f', 2, 64))
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[file.Currency]+strconv.FormatFloat(rate, 'f', 2, 64))
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[file.Currency]+amount)
	pdf.Br(24)
}

func writeTotalsWithTotalHours(pdf *gopdf.GoPdf, totalHours float64, subtotal float64, tax float64, discount float64) {
	handleRemainingSpace(pdf, (maxHeight - 600))
	pdf.SetY(600)

	writeTotal(pdf, totalHoursLabel, totalHours, false)
	writeTotal(pdf, subtotalLabel, subtotal, true)
	if tax > 0 {
		writeTotal(pdf, taxLabel, tax, true)
	}
	if discount > 0 {
		writeTotal(pdf, discountLabel, discount, true)
	}
	writeTotal(pdf, totalLabel, subtotal+tax-discount, true)
}

func writeTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64) {
	handleRemainingSpace(pdf, (maxHeight - 600))
	pdf.SetY(600)

	writeTotal(pdf, subtotalLabel, subtotal, true)
	if tax > 0 {
		writeTotal(pdf, taxLabel, tax, true)
	}
	if discount > 0 {
		writeTotal(pdf, discountLabel, discount, true)
	}
	writeTotal(pdf, totalLabel, subtotal+tax-discount, true)
}

func writeTotal(pdf *gopdf.GoPdf, label string, total float64, isMoney bool) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, label)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(12)
	pdf.SetX(amountColumnOffset - 15)
	if label == totalLabel {
		_ = pdf.SetFont("Inter-Bold", "", 11.5)
	}
  if isMoney {
	  _ = pdf.Cell(nil, currencySymbols[file.Currency]+strconv.FormatFloat(total, 'f', 2, 64))
  } else {
    _ = pdf.Cell(nil, strconv.FormatFloat(total, 'f', 2, 64))
  }
  pdf.Br(24)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
