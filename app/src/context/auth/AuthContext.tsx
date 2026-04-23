import { useState, useEffect, useCallback } from "react"
import { AuthContext } from "./useAuth"

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [username, setUsername] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch("/api/me")
      .then((r) => (r.ok ? r.json() : null))
      .then((data) => setUsername(data?.username ?? null))
      .catch((err) => {
        console.error("[AuthContext] /api/me failed:", err)
        setUsername(null)
      })
      .finally(() => setLoading(false))
  }, [])

  const login = useCallback(
    async (u: string, p: string): Promise<string | null> => {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: u, password: p }),
      })
      if (res.ok) {
        const data = await res.json()
        setUsername(data.username)
        return null
      }
      const err = await res.json()
      return err.error ?? "ログインに失敗しました"
    },
    [],
  )

  const logout = useCallback(async () => {
    await fetch("/api/logout", { method: "POST" })
    setUsername(null)
  }, [])

  return (
    <AuthContext.Provider value={{ username, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
