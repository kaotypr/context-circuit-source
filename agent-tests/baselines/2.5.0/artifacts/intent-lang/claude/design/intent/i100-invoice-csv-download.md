---
id: i100
created_by: maya
created_at: 2026-10-29T06:03:37Z
---

# Unduh invoice dalam format CSV oleh pelanggan

## Goal

Pelanggan dapat mengunduh invoice mereka sebagai file CSV langsung dari API, tanpa bantuan tim. Saat ini invoice hanya bisa diekspor sebagai plain-text summary, sehingga pelanggan yang ingin mengolah datanya di spreadsheet atau sistem akuntansi harus menyalin manual.

## Non-goals

- Mengubah atau menghapus plain-text summary yang sudah ada.
- Ekspor massal untuk admin atau tim internal.
- Mengirim CSV lewat email atau integrasi pihak ketiga.
- Mengubah aturan void dan credit note.

## Constraints

- Invoice hanya bisa di-void selama masih draft; setelah diterbitkan, koreksinya dengan credit note. CSV harus menampilkan status ini apa adanya, tidak menyembunyikannya (lihat `context/domains/invoices.md`).
- Pelanggan hanya boleh mengunduh invoice miliknya sendiri. Tidak boleh ada kebocoran data antarpelanggan.
- Kode ekspor saat ini berada di `api@internal/invoice/`; perubahan dibatasi di repository `api`.

## Success criteria

- Pelanggan yang terautentikasi dapat mengunduh CSV berisi seluruh invoice miliknya dan tidak memuat invoice pelanggan lain.
- Setiap invoice tampil dengan statusnya yang benar (draft, diterbitkan, atau dikreditkan).
- File dibuka dengan benar di spreadsheet: encoding UTF-8, kolom sesuai yang disepakati, dan angka serta tanggal dalam format yang sudah disepakati.

## Repository scope

- `api` — endpoint unduhan, penyusunan CSV, dan pengujian otorisasi. Tidak ada repository lain yang disentuh.

## Open questions

1. **Hubungan dengan plain-text summary.** CSV ditambahkan di samping plain-text summary yang sudah ada, atau menggantikannya?
2. **Granularitas baris.** Satu baris per invoice, atau satu baris per line item?
3. **Invoice yang disertakan.** Apakah invoice draft ikut diunduh pelanggan, atau hanya yang sudah diterbitkan? Dan credit note ditampilkan sebagai baris tersendiri?
4. **Kolom.** Kolom apa saja yang diinginkan? Usulan minimal: nomor invoice, tanggal terbit, jatuh tempo, status, mata uang, subtotal, pajak, total.
5. **Rentang dan ukuran.** Apakah cukup satu unduhan berisi semua invoice, atau perlu filter periode tanggal? Jika pelanggan punya banyak invoice, apakah perlu paginasi atau ekspor asinkron?
6. **Format angka dan tanggal.** Pemisah desimal dan ribuan, apakah mata uang ditulis di setiap baris, dan tanggal memakai ISO 8601 (`YYYY-MM-DD`)?
7. **Autentikasi pelanggan.** Bagaimana pelanggan terautentikasi ke API saat ini: apakah sudah ada token atau kunci API per pelanggan yang bisa dipakai untuk endpoint ini?
