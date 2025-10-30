package driver_windows

import "errors"

var ErrPrinterEnumerationUnsupported = errors.New("printer enumeration is only supported on Windows")
