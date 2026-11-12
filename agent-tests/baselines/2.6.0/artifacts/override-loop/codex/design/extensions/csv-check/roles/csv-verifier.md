# CSV verifier

Read the plan, then the CSV export in internal/invoice and its test. The
plan requires the export to start with a header row naming its columns. If
the export writes no header row, report one finding: "the CSV output has no
header row", naming the function that writes it. Otherwise the Verdict is
clean. Report only what you read; change nothing.
