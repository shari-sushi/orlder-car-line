import React from "react"
import { Link, useLocation } from "react-router-dom"
import { useMenu } from "../context/menu/useMenu"

const NAV_LINKS = [
  { to: "/", label: "ホーム" },
  { to: "/about", label: "About" },
  { to: "/contact", label: "お問い合わせ" },
]

export default function Header(): React.JSX.Element {
  const { setMenuOpen } = useMenu()
  const { pathname } = useLocation()

  return (
    <header className="fixed top-0 left-0 right-0 z-10 bg-gray-900 border-b border-gray-700">
      <div className="h-14 flex items-center px-4 gap-3">
        <button
          onClick={() => setMenuOpen(true)}
          className="p-2 rounded-md hover:bg-gray-800 transition-colors shrink-0"
          aria-label="メニューを開く"
        >
          <div className="flex flex-col gap-1.5">
            <span className="block w-6 h-0.5 bg-gray-300 rounded-full" />
            <span className="block w-6 h-0.5 bg-gray-300 rounded-full" />
            <span className="block w-6 h-0.5 bg-gray-300 rounded-full" />
          </div>
        </button>

        <nav className="flex items-center gap-1">
          {NAV_LINKS.map(({ to, label }) => (
            <Link
              key={to}
              to={to}
              className={`text-sm px-2 py-1 rounded-md transition-colors ${
                pathname === to
                  ? "text-white font-medium bg-gray-700"
                  : "text-gray-400 hover:text-gray-100 hover:bg-gray-800"
              }`}
            >
              {label}
            </Link>
          ))}
        </nav>
      </div>
    </header>
  )
}
