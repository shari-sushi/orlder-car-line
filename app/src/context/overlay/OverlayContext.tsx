import { createContext, useCallback, useState, type ReactNode } from "react"

type OverlayContextValue = {
  open: (content: ReactNode) => void
  close: () => void
}

export const OverlayContext = createContext<OverlayContextValue | null>(null)

export function OverlayProvider({
  children,
}: {
  children: ReactNode
}): React.JSX.Element {
  const [content, setContent] = useState<ReactNode | null>(null)

  const open = useCallback((c: ReactNode) => setContent(c), [])
  const close = useCallback(() => setContent(null), [])

  return (
    <OverlayContext.Provider value={{ open, close }}>
      {children}
      {content !== null && (
        <>
          {/* 半透明の背景（DOM順序で content より前 = 自然に後ろ） */}
          <div className="fixed inset-0 bg-black/40" onClick={close} />
          {/* コンテンツ（bg より後の DOM = 自然に手前） */}
          <div className="fixed inset-0 flex items-center justify-center pointer-events-none">
            <div className="pointer-events-auto">{content}</div>
          </div>
        </>
      )}
    </OverlayContext.Provider>
  )
}
