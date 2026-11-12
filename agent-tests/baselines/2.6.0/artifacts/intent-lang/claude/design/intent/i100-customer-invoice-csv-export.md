---
id: i100
created_by: maya
created_at: 2026-11-11T07:22:43Z
---

# Customer invoice CSV download

## Goal

Pelanggan dapat mengunduh invoice mereka dalam format CSV melalui API. Saat ini invoice hanya bisa diekspor sebagai ringkasan teks biasa, sehingga pelanggan yang ingin mengolah datanya di spreadsheet harus menyalin manual.

## Non-goals

- Tidak mengubah aturan void dan credit note: draft tetap bisa di-void, invoice yang sudah diterbitkan tetap dikoreksi dengan credit note.
- Tidak mengubah ringkasan teks biasa yang sudah ada. Ringkasan tetap tersedia seperti sebelumnya.
- Tidak menambah format lain selain CSV, misalnya PDF atau Excel.

## Constraints

- Pelanggan hanya dapat mengunduh invoice miliknya sendiri. Permintaan untuk invoice milik pelanggan lain ditolak.
- Perubahan hanya di repository `api`, tempat logika invoice berada (`api@internal/invoice/`).

## Success criteria

- Pelanggan dapat meminta file CSV dari invoice mereka melalui API dan menerima file yang bisa dibuka di spreadsheet.
- Isi dan bentuk CSV sesuai keputusan pada pertanyaan 1 sampai 4.
- Permintaan untuk invoice milik pelanggan lain ditolak dan tidak menghasilkan data.
- Ringkasan teks biasa tetap menghasilkan keluaran yang sama seperti sebelum perubahan.

## Repository scope

- `api` — `internal/invoice/`, tempat ringkasan teks biasa, void, dan credit note sudah diatur.

## Open questions

1. **Apakah unduhan berupa satu invoice per permintaan, atau daftar semua invoice pelanggan dengan filter tanggal?**

2. **Apa isi file CSV: satu baris per invoice (nomor, tanggal, status, total), atau satu baris per item di setiap invoice?**

3. **Apakah invoice berstatus draft dan invoice yang sudah dikoreksi dengan credit note ikut diekspor?**

4. **Konvensi CSV apa yang dipakai: pemisah kolom koma atau titik koma, serta format angka dan tanggal? Indonesia biasanya memakai koma sebagai pemisah desimal, dan itu bisa membuat kolom angka terbaca salah di spreadsheet.**
