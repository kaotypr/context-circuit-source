import { dirname, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { refreshProductKnowledge } from "./lib/product-knowledge.js"

const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
console.log(JSON.stringify(await refreshProductKnowledge(root), null, 2))
