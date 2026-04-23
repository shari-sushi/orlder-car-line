import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import path from "path"

const apiDir = path.resolve(__dirname, "api")

export default defineConfig({
  plugins: [
    react(),
    {
      name: "local-api",
      configureServer(server) {
        server.middlewares.use(async (req, res, next) => {
          const url = req.url ?? ""

          if (!url.startsWith("/api/")) {
            next()
            return
          }

          if (req.method === "OPTIONS") {
            res.setHeader("Access-Control-Allow-Origin", "*")
            res.statusCode = 204
            res.end()
            return
          }

          const pathname = url.split("?")[0]
          const apiPath = pathname.replace("/api/", "")

          try {
            const filePath = path.resolve(apiDir, `${apiPath}.ts`)
            const mod = await import(filePath)
            const { adaptVercel } = await import(
              path.resolve(apiDir, "_lib/vercel-adapter.ts")
            )
            res.setHeader("Content-Type", "application/json")
            await adaptVercel(mod.default)(req, res)
          } catch (e) {
            console.error("[local-api] error:", e)
            res.statusCode = 500
            res.end(JSON.stringify({ error: String(e) }))
          }
        })
      },
    },
  ],
  server: {
    fs: {
      allow: [path.resolve(__dirname), path.resolve(__dirname, "..")],
    },
  },
})
