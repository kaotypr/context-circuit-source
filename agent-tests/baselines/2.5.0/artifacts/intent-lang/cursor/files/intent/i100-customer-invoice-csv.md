---
id: i100
created_by: maya
created_at: 2026-10-29T06:05:18Z
---

# Pelanggan mengunduh invoice sebagai CSV

## Goal

Pelanggan bisa mengunduh invoice miliknya sendiri dari API sebagai file CSV.

## Non-goals

- Mengubah kapan invoice di-void dan kapan dipakai credit note.
- Menyediakan unduhan dari antarmuka lain selain API.
- Menambah format lain selain CSV.

## Constraints

- Catatan [Invoices](../context/domains/invoices.md) menetapkan bahwa invoice hanya diekspor sebagai ringkasan teks biasa. CSV menyentuh aturan itu, jadi pertanyaan 1 harus dijawab sebelum ada yang membangunnya.
- Yang diunduh hanya invoice milik pelanggan yang meminta.

## Success criteria

- Pelanggan yang sudah terautentikasi menerima file CSV berisi invoice miliknya lewat API.
- Invoice milik pelanggan lain tidak ikut dalam file itu.
- Draft tetap hanya bisa di-void, dan invoice yang sudah terbit tetap memakai credit note.

## Repository scope

- `api` — ekspor invoice, termasuk ringkasan teks yang ada sekarang.

## Open questions

1. **Apakah CSV menggantikan ringkasan teks biasa, atau berdiri di sampingnya sehingga pelanggan tetap bisa mengunduh keduanya?**
