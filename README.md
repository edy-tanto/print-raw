# Print Raw Web Service

Service Windows untuk mencetak struk POS dan dokumen terkait via HTTP API.

## Prasyarat

- **Windows** 7 atau lebih baru
- **Go 1.17.x** (untuk build dari sumber)
- Printer receipt terpasang di Windows atau terjangkau via jaringan
- **Git Bash** (untuk menjalankan script build)

## Build

Build menghasilkan satu paket ZIP berisi executable 32-bit dan aset yang dibutuhkan.

1. Buka terminal di **root project** (Git Bash di Windows).

2. Jalankan script build:
   ```bash
   sh build.sh
   ```

3. Saat diminta, masukkan **version** hasil build (contoh: `v0.5` atau `0.5`).  
   - Kosongkan lalu Enter = pakai version yang ada di `bin/version.txt`.

4. Hasil build:
   - Binary: `bin/print_web_service.exe`
   - ZIP installer: `installer-printer-<version>.zip` (disimpan ke `C:\Users\Dream\Downloads`)

Detail build: Go 32-bit (`GOARCH=386`), executable + file `.bmp` dan script di `bin/` ikut dibundle ke ZIP.

## Menjalankan

### Opsi 1: Development (go run)

Jalankan tanpa build, dari **root project**:

```bash
go run ./cmd/print_web_service
```

Atau di PowerShell:

```powershell
go run .\cmd\print_web_service
```

Server HTTP berjalan di **http://localhost:8080**. Hentikan dengan Ctrl+C.  
Pastikan aset `.bmp` yang dipakai ada di `bin/` (atau root project, tergantung resolusi path di kode).

### Opsi 2: Langsung (executable)

Jalankan binary yang sudah di-build dari folder `bin/`:

```powershell
cd bin
.\print_web_service.exe
```

Server HTTP berjalan di **http://localhost:8080**.  
Pastikan folder `bin/` berisi file `.bmp` yang dipakai (header receipt, dll).

### Opsi 3: Sebagai Windows Service

Untuk dipakai di production (jalan di background, auto-start):

1. **Install service** (jalankan sebagai Administrator):
   ```powershell
   cd bin
   .\install.bat
   ```

2. **Kontrol service:**
   - Start: `bin\start.bat`
   - Stop: `bin\stop.bat`
   - Uninstall: `bin\uninstall.bat`

3. Cek status:
   ```powershell
   sc.exe query PrintRawWeb
   ```

Dokumentasi API dan contoh payload: [docs/print_raw_web_service.md](docs/print_raw_web_service.md).
