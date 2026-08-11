import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

export async function git(cwd: string, args: string[]): Promise<string> {
  try {
    const { stdout } = await execFileAsync("git", args, {
      cwd,
      encoding: "utf8",
      maxBuffer: 10 * 1024 * 1024,
    });
    return stdout.trim();
  } catch (error) {
    const detail = error as Error & { stderr?: string };
    throw new Error(`git ${args.join(" ")} failed in ${cwd}: ${detail.stderr?.trim() || detail.message}`);
  }
}

export async function assertCleanRepository(repository: string): Promise<void> {
  await git(repository, ["rev-parse", "--is-inside-work-tree"]);
  const dirty = await git(repository, ["status", "--porcelain=v1", "--untracked-files=normal"]);
  if (dirty) {
    throw new Error(`Repository has unresolved local changes; refusing to continue: ${repository}\n${dirty}`);
  }
}
