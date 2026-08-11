import type { ResetCounterResponse } from "./contract.js";

export function resetCounter(): ResetCounterResponse {
  return { state: { count: 0 } };
}
