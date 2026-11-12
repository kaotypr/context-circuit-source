---
id: i100
created_by: maya
created_at: 2026-11-11T16:38:38Z
---

# Customer invoice CSV download

## Goal

Customer bisa download invoice milik mereka sendiri sebagai file CSV, langsung dari API.

## Non-goals

- Ringkasan plain-text yang sekarang ada tetap jalan. CSV nambah format, nggak nggantiin ringkasan itu.
- Nggak ada halaman atau button baru. Download-nya lewat API.
- Aturan void dan credit note nggak berubah. Draft tetap di-void, invoice yang sudah terbit tetap di-credit.

## Constraints

- Satu customer cuma bisa ambil invoice miliknya sendiri.
- Export yang ada sekarang cuma ringkasan plain-text, jadi CSV ini jalur baru di API.

## Success criteria

- Customer yang sudah login bisa minta CSV invoice mereka dari API, dan responsenya file CSV yang bisa di-download.
- File itu cuma berisi invoice milik customer yang minta.
- Ringkasan plain-text masih bisa diambil seperti sekarang.

## Repository scope

- `api` — endpoint download CSV dan isinya.

## Open questions

1. **Satu baris CSV itu satu invoice, atau satu line item?**
   _Answer:_
2. **Draft ikut di file, atau cuma invoice yang sudah terbit?**
   _Answer:_
