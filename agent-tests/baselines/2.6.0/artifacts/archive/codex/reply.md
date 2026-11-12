p1000 is marked complete. I updated the invoice knowledge to include CSV exports, and `check` returned `ok: true` with no issues.

The CLI also moved the archived i101 and p1001 records into `_archive/`. These changes are uncommitted. Since this was the first changing command on CLI 2.5.0, one person should commit and push the migration before anyone else changes the workspace; other members should install 2.5.0 or later before pulling.

Would you like me to **commit and push** these files, or **commit only**?