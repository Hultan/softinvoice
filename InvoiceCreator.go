package main

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
)

type InvoiceCreator struct {
	Invoice *Invoice
	PixBuf  *gdk.Pixbuf
	Surface *cairo.Surface
	Context *cairo.Context
}
