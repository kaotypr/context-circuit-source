const secretPattern = /(?:token|password|passwd|secret|api[_-]?key)\s*[=:]|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/i;
const schemePattern = /^[A-Za-z][A-Za-z0-9+.-]*:\/\//;
const scpPattern = /^(?:git|ssh)@[A-Za-z0-9.-]+:[A-Za-z0-9._~/-]+$/;
const providerPattern = /^(?:github|gitlab|linear|jira|notion):[A-Za-z0-9][A-Za-z0-9._/@#:+-]*$/;
const localSourcePattern = /^(?!.*(?:^|\/)\.\.(?:\/|$))[A-Za-z0-9._@#:+%=-]+(?:\/[A-Za-z0-9._@#:+%=-]+)*$/;

function decoded(value: string): string {
  try { return decodeURIComponent(value); } catch { return value; }
}

function urlError(value: string, allowedSchemes: Set<string>): string | null {
  let parsed: URL;
  try { parsed = new URL(value); } catch { return "is not a valid URL"; }
  if (!allowedSchemes.has(parsed.protocol)) return `uses unsupported URL scheme ${parsed.protocol}`;
  const authority = value.slice(value.indexOf("://") + 3).split(/[/?#]/, 1)[0] ?? "";
  if (authority.includes("@") || parsed.username || parsed.password || decoded(parsed.username) || decoded(parsed.password)) return "appears to contain credentials in URL userinfo";
  if (secretPattern.test(decoded(value))) return "appears to contain credentials";
  return null;
}

export function remoteReferenceError(value: string): string | null {
  const reference = value.trim();
  if (!reference || /[\r\n]/.test(reference)) return "must be a non-empty single-line reference";
  if (secretPattern.test(decoded(reference))) return "appears to contain credentials";
  if (schemePattern.test(reference)) return urlError(reference, new Set(["https:", "ssh:", "git:"]));
  if (scpPattern.test(reference)) return null;
  if (/^[^\s@/:]+@[^\s:]+:/.test(reference)) return "contains unsupported remote userinfo";
  if (/https?:/i.test(reference) || /%40/i.test(reference)) return "contains invalid or encoded URL userinfo";
  return "must be a credential-free HTTPS, SSH, Git, or SCP-style remote";
}

export function cloneReferenceError(value: string): string | null {
  const remoteError = remoteReferenceError(value);
  if (!remoteError) return null;
  const reference = value.trim();
  if (secretPattern.test(decoded(reference)) || /[\r\n]/.test(reference) || /%40/i.test(reference)) return remoteError;
  if (reference.startsWith("/") || localSourcePattern.test(reference)) return null;
  return remoteError;
}

export function contextReferenceError(value: string): string | null {
  const reference = value.trim();
  if (!reference || /[\r\n\\]/.test(reference)) return "must be a non-empty single-line reference without backslashes";
  const decodedReference = decoded(reference);
  if (secretPattern.test(decodedReference)) return "appears to contain credentials";
  if (/\\/.test(decodedReference)) return "must not contain encoded backslashes";
  if (schemePattern.test(reference)) return urlError(reference, new Set(["https:"]));
  if (providerPattern.test(reference)) return null;
  if (reference.startsWith("/") || decodedReference.startsWith("/") || /^[A-Za-z]:/.test(decodedReference)) return "must not be an absolute local path";
  if (decodedReference.split("/").some((part) => part === ".." || part === ".")) return "must not contain traversal segments";
  if (!localSourcePattern.test(reference)) return "must be a safe relative path or allowed provider reference";
  return null;
}
