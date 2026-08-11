import { access, cp, mkdir } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { git } from "./lib/git.js";

const workspaceRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const source = join(workspaceRoot, "fixtures", "react-app");
const destination = resolve(workspaceRoot, "repositories", "frontend");

await mkdir(join(workspaceRoot, "repositories"), { recursive: true });
try {
  await access(destination);
  throw new Error(`Fixture destination already exists; preserving it: ${destination}`);
} catch (error) {
  if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
}
await cp(source, destination, { recursive: true, errorOnExist: true, force: false });
await git(destination, ["init", "--initial-branch=main"]);
await git(destination, ["add", "."]);
await git(destination, ["-c", "user.name=Kao Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "chore: initialize React fixture"]);
console.log(`Created fixture repository: ${destination}`);
