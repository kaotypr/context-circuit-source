---
id: i100
created_by: maya
created_at: 2026-11-12T09:41:40Z
---

# Unduh invoice sebagai CSV melalui API

## Goal

Customer dapat mengunduh invoice miliknya sebagai file CSV melalui API. Informasi invoice yang sudah tersedia tetap bisa diakses seperti sebelumnya.

## Non-goals

Perubahan ini tidak mengubah cara invoice dibuat, diterbitkan, dibatalkan, atau dikreditkan.

## Constraints

Invoice yang sudah diterbitkan hanya dapat dikreditkan dengan credit note; aturan siklus invoice ini tetap berlaku. Saat ini API menyediakan ringkasan invoice sebagai plain text.

## Success criteria

- Customer dapat meminta invoice miliknya melalui API dan menerima file CSV yang dapat dibuka sebagai data tabular.

## Repository scope

`api` — menyediakan akses invoice dan respons unduhan CSV.

## Open questions

1. **Informasi apa saja yang perlu ada di CSV?**
   _Answer:_ Belum diputuskan. Rekomendasinya mencakup informasi yang saat ini ada di ringkasan invoice, agar CSV tetap berguna tanpa memperluas isi invoice secara diam-diam.
