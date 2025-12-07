package utils

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatDatetime(dateString string) string {
	layout := "2006-01-02T15:04:05.000"
	parsedTime, _ := time.ParseInLocation(layout, dateString, time.Local)
	localTime := parsedTime.In(time.Local)
	return localTime.Format("02/01/2006 15:04:05")
}

func FormatMoney(amount float32) string {
	p := message.NewPrinter(language.MustParse("id"))
	return p.Sprintf("%.0f", amount)
}

func FormatMoneyTwoDigitAfterComma(amount float32) string {
	p := message.NewPrinter(language.MustParse("id"))
	return p.Sprintf("%.2f", amount)
}

func CenterInParentheses(text string, width int) string {
	// Account for the two parentheses characters
	availableWidth := width - 2
	if len(text) > availableWidth {
		// Truncate text if it's too long to fit
		text = text[:availableWidth]
	}
	padding := availableWidth - len(text)
	left := padding / 2
	right := padding - left
	return fmt.Sprintf("(%s%s%s)", strings.Repeat(" ", left), text, strings.Repeat(" ", right))
}

func CenterText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := width - len(text)
	left := padding / 2
	right := padding - left
	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", left), text, strings.Repeat(" ", right))
}
