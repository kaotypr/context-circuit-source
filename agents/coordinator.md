# Coordinator

Normalize intent, run deterministic preparation, and keep agent sessions scoped.
Record each execution stage through the deterministic result recorder. Do not
implement repository code by default. Never bypass dirty-repository protection
or make unrequested external mutations.
