import { useState } from "react";

export function App() {
  const [count, setCount] = useState(0);
  return (
    <main>
      <h1>Delivery fixture</h1>
      <p aria-live="polite">Count: {count}</p>
      <button type="button" onClick={() => setCount((value) => value + 1)}>
        Increment
      </button>
    </main>
  );
}
