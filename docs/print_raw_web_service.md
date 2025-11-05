# PrintRawWeb POS Printing Service

This document explains how to build, install, and call the `PrintRawWeb` Windows service to print POS receipts and related documents.

## Overview

- **Binary:** `print_web_service.exe`
- **Default port:** `8080`
- **Primary endpoint:** `POST /` (prints a POS sales receipt by default)
- **Supporting endpoints:** `/printers`, `/cash-refund`, `/kitchen`, `/kitchen-eth`, `/table-check`, `/captain-order-bill`, `/captain-order-invoice`, `/shift/report/cash-count`

The service accepts JSON payloads describing receipts, converts any referenced bitmaps from the executable directory into ESC/POS-compatible bytes, and sends them to the configured printer.

## Prerequisites

- Windows 7 or later with Administrator access (required to register services).
- Go toolchain (1.17.x) if you need to rebuild the binary.
- Target receipt printer installed locally or reachable over Ethernet.
- Bitmap assets (e.g., `paradis-q.bmp`, `captain-order-receipt-header.bmp`) copied into the same directory as the executable.

## Build & Deploy

1. **Build the service binary (optional if you already have it):**
   ```powershell
   go build -o bin\print_web_service.exe .\cmd\print_web_service
   ```
2. **Stage runtime assets:**
   - Copy `bin\print_web_service.exe` and required `.bmp` header images into the folder where the service will run (the install scripts expect everything beside each other in `bin\`).
3. **Install or update the Windows service (run as Administrator):**
   ```powershell
   bin\install.bat
   ```
   - If the service already exists, the script simply starts it.
   - Use `bin\uninstall.bat` to remove and `bin\start.bat` / `bin\stop.bat` to control it afterward.
4. **Verify status (optional):**
   ```powershell
   sc.exe query PrintRawWeb
   ```

> **Note:** When running as a service, Windows defaults the working directory to `C:\Windows\System32`. The application automatically resolves relative asset paths against the executable location, so keep any bitmap assets alongside `print_web_service.exe`.

## API Summary

| Method | Path                     | Description                                    |
|--------|--------------------------|------------------------------------------------|
| GET    | `/printers`              | Lists available local printers.                |
| POST   | `/`                      | Prints a POS sales receipt.                    |
| POST   | `/cash-refund`           | Prints waterpark cash refund slips.            |
| POST   | `/kitchen`               | Prints kitchen orders for Patio & Dimsum.      |
| POST   | `/kitchen-eth`           | Sends kitchen orders over Ethernet printers.   |
| POST   | `/table-check`           | Prints table check dockets.                    |
| POST   | `/captain-order-bill`    | Prints captain order bills.                    |
| POST   | `/captain-order-invoice` | Prints captain order invoices.                 |
| POST   | `/shift/report/cash-count` | Prints shift cash-count reports.             |

All POST endpoints expect JSON matching the DTOs in `internal/print_web_service/dto`. Unrecognized methods return `404`, and each handler supports CORS preflight with `OPTIONS`.

## POS Receipt Payload

`POST /` expects the following shape (`dto.PrintRequestBody`):

```json
{
  "sales": {
    "id": 12345,
    "unit_business_name": "Paradis Bistro",
    "code": "INV-2024-0001",
    "op": "OP-78311",
    "customer_name": "Jane Doe",
    "table_number": "A7",
    "payment_method": "Credit Card",
    "date": "2024-05-10 19:32",
    "is_print_as_copy": false,
    "footnote": "Thank you for dining with us!",
    "footnote_align": "center",
    "grand_total": 485000,
    "credit_card_charge": 1500,
    "sales_details": [
      {
        "item": "Grilled Salmon",
        "qty": 2,
        "total_final": 240000,
        "subtotal_with_tax": 240000
      },
      {
        "item": "Sparkling Water",
        "qty": 1,
        "total_final": 45000,
        "subtotal_with_tax": 45000
      }
    ]
  },
  "printer_name": "EPSON TM-T88V Receipt"
}
```

Numbers are interpreted as floats/ints by Go, so you can send either integers or decimals. `printer_name` must exactly match a Windows-installed printer or the Ethernet target in other DTOs.

## cURL Examples

List printers:

```bash
curl --request GET \
     --url http://localhost:8080/printers
```

Print a POS sales receipt:

```bash
curl --request POST \
     --url http://localhost:8080/ \
     --header "Content-Type: application/json" \
     --data @- <<'JSON'
{
  "sales": {
    "id": 12345,
    "unit_business_name": "Paradis Bistro",
    "code": "INV-2024-0001",
    "op": "OP-78311",
    "customer_name": "Jane Doe",
    "table_number": "A7",
    "payment_method": "Credit Card",
    "date": "2024-05-10 19:32",
    "is_print_as_copy": false,
    "footnote": "Thank you for dining with us!",
    "footnote_align": "center",
    "grand_total": 485000,
    "credit_card_charge": 1500,
    "sales_details": [
      {
        "item": "Grilled Salmon",
        "qty": 2,
        "total_final": 240000,
        "subtotal_with_tax": 240000
      },
      {
        "item": "Sparkling Water",
        "qty": 1,
        "total_final": 45000,
        "subtotal_with_tax": 45000
      }
    ]
  },
  "printer_name": "EPSON TM-T88V Receipt"
}
JSON
```

The service echoes the submitted payload if the print job is enqueued successfully.

## Postman Collection

Import `docs/postman/PrintRawWeb.postman_collection.json` into Postman. It contains ready-to-run requests for listing printers and triggering POS receipt printing. Update the `baseUrl` variable if you deploy the service to a different host or port.
