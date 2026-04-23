import { createContext, useContext } from "react"

export interface AuthState {
  username: string | null
  loading: boolean
  login: (username: string, password: string) => Promise<string | null>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthState>({
  username: null,
  loading: true,
  login: async () => null,
  logout: async () => {},
})

export function useAuth() {
  return useContext(AuthContext)
}
