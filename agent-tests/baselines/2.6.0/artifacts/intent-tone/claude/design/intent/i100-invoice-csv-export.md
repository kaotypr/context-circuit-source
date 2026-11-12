---
id: i100
created_by: maya
created_at: 2026-11-11T07:23:30Z
---

# Export invoice ke CSV lewat API

## Goal

Customer bisa download invoice-invoice mereka sebagai file CSV langsung dari API. Saat ini invoice cuma bisa diexport sebagai ringkasan teks biasa, jadi customer yang mau olah datanya di spreadsheet atau aplikasi akuntansi harus ngetik ulang. Setelah ini, satu request ke API ngasih CSV berisi invoice milik customer itu, dengan kolom yang bisa langsung dibuka di spreadsheet.

## Non-goals

- Tombol download di web UI. Yang dikerjain cuma endpoint di API.
- Ngubah ringkasan teks yang sudah ada. Ringkasan teks tetap jalan seperti biasa.
- Export PDF, export terjadwal, atau kirim CSV lewat email.
- Ngubah cara invoice dibuat, diterbitkan, di-void, atau di-credit.

## Constraints

- Customer cuma dapat invoice miliknya sendiri. Tidak ada data customer lain yang ikut keluar.
- Nominal ditulis dalam sen (minor units), sama seperti `Line.Cents` dan ringkasan teks yang sudah ada, supaya angkanya nggak beda antar format.
- Invoice yang sudah diterbitkan nggak pernah di-void, cuma di-credit. Jadi CSV harus bisa nunjukin credit note, bukan cuma invoice asalnya.

## Success criteria

- Customer yang punya invoice bisa minta CSV dan dapat semua invoice miliknya dalam rentang yang dipilih, dengan baris header.
- Total per invoice di CSV sama dengan total yang muncul di ringkasan teks untuk invoice yang sama.
- Customer lain tidak bisa mengambil invoice orang lain lewat endpoint ini.
- Export tetap selesai untuk akun dengan banyak invoice, tanpa timeout, dalam rentang yang dibatasi.

## Repository scope

- `api`: endpoint export, format CSV, dan pengecekan akses. Kode invoice yang ada ada di `internal/invoice/`.

## Open questions

1. **Satu baris per invoice atau per line item?** Kalau per line item, invoice yang punya banyak baris muncul berulang dengan nomor invoice yang sama, dan total invoice nggak bisa dijumlah langsung dari CSV. Kalau per invoice, file lebih ringkas tapi detail barang nggak kelihatan. Rekomendasinya per line item, karena customer yang minta CSV biasanya mau detailnya buat pembukuan, dan total per invoice tetap bisa dihitung dengan menjumlahkan barisnya.
2. **Rentang tanggal dan batasnya?** Customer bisa milih rentang tanggal penerbitan, dan satu request dibatasi maksimal dua belas bulan supaya akun yang invoice-nya banyak nggak timeout. Batas ini bisa dinaikin nanti kalau perlu. Setuju dengan batas dua belas bulan, atau mau angka lain?
3. **Draft dan credit note ikut masuk nggak?** Rekomendasinya draft nggak ikut, karena belum pernah diterbitkan dan customer belum seharusnya melihatnya. Credit note masuk sebagai baris sendiri, dengan kolom yang nunjukin invoice mana yang di-credit, supaya saldo yang keluar dari CSV sesuai dengan yang tampil di akun.
4. **Siapa yang boleh export?** Rekomendasinya siapa pun yang punya akses lihat ke invoice di akun itu, bukan cuma owner-nya. Kalau ada role yang nggak boleh lihat nominal, role itu juga nggak boleh export. Setuju, atau cuma owner yang boleh?