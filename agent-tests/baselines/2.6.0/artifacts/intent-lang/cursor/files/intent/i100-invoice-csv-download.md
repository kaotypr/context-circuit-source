---
id: i100
created_by: maya
created_at: 2026-11-11T16:36:49Z
---

# Customer invoice CSV download

## Goal

Customer yang sudah bisa melihat invoice-nya bisa mengunduh invoice itu sebagai CSV lewat API. Satu permintaan menghasilkan satu file, satu baris per invoice, dan isi barisnya mengikuti ringkasan teks yang sekarang sudah dihasilkan untuk invoice itu.

## Non-goals

Aturan void dan credit note tidak berubah: draft tetap di-void, invoice yang sudah terbit tetap dikredit lewat credit note. Tidak ada halaman unduh di luar API, dan tidak ada format lain selain CSV.

## Constraints

Hari ini invoice hanya diekspor sebagai ringkasan teks. CSV ini menambah cara unduh, dan barisnya memakai fakta yang sama dengan ringkasan itu, termasuk status void atau credit kalau invoice itu sudah lewat jalur tersebut.

## Success criteria

- Customer yang berhak melihat suatu invoice bisa meminta CSV-nya dari API dan membuka file itu di spreadsheet.
- Setiap baris mewakili satu invoice, dan kolomnya memuat isi yang sama dengan ringkasan teks invoice itu.
- Invoice yang customer itu tidak boleh lihat tidak ikut masuk ke file.
- Meminta ringkasan teks tetap berhasil seperti sekarang, selama pertanyaan 1 dijawab dengan tetap menyimpannya.

## Repository scope

Perubahan ada di repositori `api`, pada bagian yang sekarang mengurus void, credit note, dan ringkasan teks invoice.

## Open questions

1. **Ringkasan teks tetap tersedia setelah CSV ada, atau CSV menggantikannya?**
2. **Satu baris per invoice, atau satu baris per line item?**
3. **Invoice yang masih draft ikut masuk ke file, atau hanya invoice yang sudah terbit?**
4. **Semua customer yang sudah bisa melihat invoice itu boleh mengunduh, atau hanya pemilik akun?**
