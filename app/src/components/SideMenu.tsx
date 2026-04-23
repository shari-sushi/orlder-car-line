import React from "react"
import { useNavigate, useLocation } from "react-router-dom"
import { useMenu } from "../context/menu/useMenu"

const NAV_ITEMS = [
  { to: "/", label: "ホーム" },
  { to: "/about", label: "About" },
  { to: "/contact", label: "お問い合わせ" },
]

export default function SideMenu(): React.JSX.Element {
  const { menuOpen, setMenuOpen } = useMenu()
  const navigate = useNavigate()
  const { pathname } = useLocation()

  const handleNav = (to: string) => {
    setMenuOpen(false)
    navigate(to)
  }

  return (
    <>
      {/* オーバーレイ */}
      <div
        className={`fixed inset-0 bg-black/40 transition-opacity duration-300 ${menuOpen ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"}`}
        onClick={() => setMenuOpen(false)}
      />

      {/* パネル */}
      <div
        className={`fixed top-0 left-0 bottom-0 bg-gray-900 shadow-2xl flex flex-col w-full md:w-80 transition-transform duration-300 ${menuOpen ? "translate-x-0" : "-translate-x-full"}`}
      >
        {/* ヘッダー */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-gray-700">
          <h2 className="text-lg font-bold text-gray-100">メニュー</h2>
          <button
            onClick={() => setMenuOpen(false)}
            className="w-8 h-8 flex items-center justify-center rounded-md hover:bg-gray-800 transition-colors text-gray-400 text-xl"
            aria-label="閉じる"
          >
            ✕
          </button>
        </div>

        {/* ナビ */}
        <nav className="flex-1 overflow-y-auto p-4">
          <ul className="space-y-1">
            {NAV_ITEMS.map(({ to, label }) => (
              <li key={to}>
                <button
                  onClick={() => handleNav(to)}
                  className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    pathname === to
                      ? "bg-blue-900 text-blue-300"
                      : "text-gray-300 hover:bg-gray-800"
                  }`}
                >
                  {label}
                </button>
              </li>
            ))}
          </ul>
        </nav>
      </div>
    </>
  )
}
