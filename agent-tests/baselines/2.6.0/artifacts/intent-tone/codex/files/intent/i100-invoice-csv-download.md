---
id: i100
created_by: maya
created_at: 2026-11-12T09:42:27Z
---

# Download invoice sebagai CSV

## Goal

Customer bisa mengunduh invoice mereka dalam format CSV lewat API, supaya data invoice mudah dipakai di spreadsheet atau aplikasi akuntansi.

## Non-goals

Perubahan ini tidak mengubah cara customer membuat atau mengirim invoice, dan tidak mengganti plain-text summary yang sudah tersedia.

## Constraints

Invoice saat ini hanya bisa diekspor sebagai plain-text summary. CSV menjadi format export tambahan.

## Success criteria

- Customer bisa meminta dan mengunduh CSV berisi invoice mereka lewat API.
- File yang diunduh bisa dibuka sebagai CSV dan memuat data invoice sesuai pilihan yang disepakati.
- Akses ke invoice tetap mengikuti batas akses customer yang berlaku.

## Repository scope

- `api` — menyediakan download invoice dalam format CSV.

## Open questions

1. **Invoice mana yang masuk ke file CSV?**
2. **Bagaimana detail invoice disusun dalam baris dan kolom CSV?**
