---
id: i100
created_by: maya
created_at: 2026-10-29T06:04:06Z
---

# Customers can download invoices as CSV

## Goal

Pelanggan dapat mengunduh data invoice mereka melalui API dalam format CSV, sehingga mereka bisa memakai data tersebut di spreadsheet atau alur kerja lain.

## Non-goals

Perubahan pada aturan invoice, format ekspor selain CSV, atau perubahan pada antarmuka non-API.

## Constraints

Ekspor yang ada saat ini berupa plain-text summary. Perubahan ini harus menentukan bagaimana ekspor CSV berhubungan dengan format yang sudah ada tanpa mengubah aturan invoice.

## Success criteria

- Pelanggan dapat meminta unduhan CSV invoice melalui API dan menerima file CSV yang dapat dibaca oleh spreadsheet.
- CSV berisi data invoice yang diminta dan menggunakan nama kolom serta format nilai yang konsisten.
- Akses unduhan mengikuti otorisasi pelanggan atas invoice.

## Repository scope

`api` — endpoint API dan pembuatan respons unduhan CSV.

## Open questions

1. **Apakah CSV harus menggantikan plain-text summary yang ada, atau tersedia sebagai format unduhan tambahan?**

   _Answer:_
