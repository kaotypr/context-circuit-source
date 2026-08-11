import assert from "node:assert/strict";
import test from "node:test";
import { resetCounter } from "./api.js";

test("resetCounter returns the contract reset state", () => {
  assert.deepEqual(resetCounter(), { state: { count: 0 } });
});
