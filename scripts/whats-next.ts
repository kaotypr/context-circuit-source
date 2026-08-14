import { dirname, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { recommendWhatsNext } from "./lib/whats-next.js"

const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
console.log(JSON.stringify(await recommendWhatsNext(root), null, 2))
