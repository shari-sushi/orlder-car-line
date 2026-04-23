import { createContext, useState, type ReactNode } from "react"

interface MenuContextValue {
  menuOpen: boolean
  setMenuOpen: (open: boolean) => void
}

// eslint-disable-next-line react-refresh/only-export-components -- context object must be exported for useMenu hook
export const MenuContext = createContext<MenuContextValue | null>(null)

export function MenuProvider({
  children,
}: {
  children: ReactNode
}): React.JSX.Element {
  const [menuOpen, setMenuOpen] = useState(false)
  return (
    <MenuContext.Provider value={{ menuOpen, setMenuOpen }}>
      {children}
    </MenuContext.Provider>
  )
}
