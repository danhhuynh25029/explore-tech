package main

type OldPrinter interface {
	PrintOldMessage() string
}

type LegacyPrinter struct{}

func (lp *LegacyPrinter) PrintOldMessage() string {
	return "Legacy Printer: Old message"
}

type PrinterAdapter struct {
	oldPrinter *LegacyPrinter
}

func NewPrinterAdapter(oldPrinter *LegacyPrinter) *PrinterAdapter {
	return &PrinterAdapter{oldPrinter}
}

func (npa *PrinterAdapter) PrintOldMessage() string {
	return npa.oldPrinter.PrintOldMessage() + " - adapted"
}

func main() {
	var printers []OldPrinter
	printers = append(printers, &LegacyPrinter{})
	adapter := NewPrinterAdapter(&LegacyPrinter{})
	printers = append(printers, adapter)

	for _, v := range printers {
		v.PrintOldMessage()
	}
}
