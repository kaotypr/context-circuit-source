import { readFile } from "node:fs/promises"
import { dirname, resolve } from "node:path"
import { parseArgs } from "node:util"
import { fileURLToPath } from "node:url"
import { importProductKnowledge } from "./lib/product-knowledge.js"
import type { ProductKnowledgeImportRequest } from "./lib/types.js"

const { values } = parseArgs({ options: { source: { type: "string" }, request: { type: "string" }, "source-id": { type: "string" }, kind: { type: "string" }, page: { type: "string", multiple: true, default: [] }, title: { type: "string" }, purpose: { type: "string" } } })
if (!values.source && !values.request) throw new Error("Usage: cc import-product-knowledge --source <path> | --request <json-file>")
const root = resolve(process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT ?? resolve(dirname(fileURLToPath(import.meta.url)), ".."))
const request = values.request ? JSON.parse(await readFile(resolve(process.cwd(), values.request), "utf8")) as ProductKnowledgeImportRequest : { source: values.source!, ...(values["source-id"] ? { source_id: values["source-id"] } : {}), ...(values.kind ? { kind: values.kind } : {}), ...(values.page?.length ? { product_knowledge: values.page } : {}), ...(values.title ? { title: values.title } : {}), ...(values.purpose ? { purpose: values.purpose } : {}) }
console.log(JSON.stringify(await importProductKnowledge(root, request), null, 2))
