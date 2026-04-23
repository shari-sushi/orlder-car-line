import type { ReactNode } from "react"
import { useOverlay } from "../context/overlay/useOverlay"

interface ModalProps {
  children: ReactNode
}

/**
 * useOverlay().open(<Modal>...</Modal>) で使うカードコンポーネント。
 * 外側のオーバーレイ背景は OverlayContext が担当する。
 */
export default function Modal({ children }: ModalProps): React.JSX.Element {
  const { close } = useOverlay()

  return (
    <div className="bg-white rounded-xl shadow-xl p-6 w-80 relative">
      <button
        onClick={close}
        className="absolute top-3 right-3 w-7 h-7 flex items-center justify-center rounded-md hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
        aria-label="閉じる"
      >
        ✕
      </button>
      {children}
    </div>
  )
}
