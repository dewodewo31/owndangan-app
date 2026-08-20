// jest jsdom needs the browser-stdlib encoder/decoder that react-dom's
// server HTML build references. Node provides them via `util`.
const { TextEncoder, TextDecoder } = require("util")
Object.assign(globalThis, { TextEncoder, TextDecoder })

// framer-motion console.errors a "useLayoutEffect does nothing on the server"
// warning per motion div while react-dom renders to static markup; that noise
// drowns real assertion failures, so drop only that specific warning.
const realError = console.error
console.error = (...args) => {
  const first = typeof args[0] === "string" ? args[0] : ""
  if (first.startsWith("Warning: useLayoutEffect does nothing on the server")) return
  realError(...args)
}