export function lsGet<T>(key: string): T | null {
  try {
    const stored = localStorage.getItem(key)
    if (stored) return JSON.parse(stored) as T
  } catch (err) {
    console.error(`[localStorage] parse failed (key: ${key}):`, err)
  }
  return null
}

export function lsSet(key: string, value: unknown): void {
  localStorage.setItem(key, JSON.stringify(value))
}
