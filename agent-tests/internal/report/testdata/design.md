## Before / after

Each pair is one scenario and one host. The run on a side is the passing run with the median total input tokens. The lowest run number wins a tie. When no run passed, the median of the remaining runs is used.

### `intent` · claude

Pair: branch r1 (1 run), before r1 (1 run). Chosen as the passing run with the median total input tokens.

| | Before | After |
|---|---|---|
| Files produced | `design/intent/i100-demo.md`, `files/note.txt` | `design/intent/i100-demo.md`, `files/note.txt` |
| Headings | Request | Request · Open questions |
| Lines | 6 | 10 |

<details><summary>design/intent/i100-demo.md (before → after)</summary>

```diff
--- before/design/intent/i100-demo.md
+++ after/design/intent/i100-demo.md
@@ -1,3 +1,7 @@
 # Request
 
 Customers need a CSV download.
+
+# Open questions
+
+None known.
```
</details>

<details><summary>files/note.txt (before → after)</summary>

```diff
--- before/files/note.txt
+++ after/files/note.txt
@@ -1,3 +1,3 @@
 alpha
 beta
-gamma
+delta
```
</details>

### Reads

```diff
--- before/reads.txt
+++ after/reads.txt
@@ -1,1 +1,2 @@
 context/INDEX.md
+context/domains/invoices.md
```

### Reply excerpts

First 12 lines. Key-shaped strings are replaced with `[redacted]`.

#### Before

```
The before reply.
It asks for approval.
```

#### After

```
The after reply.
It still asks for approval.
```

### Tokens

Total input tokens: 198k → 212k (+7%)
Output tokens: 3.9k → 4.1k (+5%)
