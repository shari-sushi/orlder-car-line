import { BrowserRouter, Routes, Route, useLocation } from "react-router-dom"
import { useEffect } from "react"
import TopPage from "./pages/TopPage"
import MdPage from "./pages/MdPage"
import Header from "./components/Header"
import SideMenu from "./components/SideMenu"
import { MenuProvider } from "./context/menu/MenuContext"
import { OverlayProvider } from "./context/overlay/OverlayContext"
import { AuthProvider } from "./context/auth/AuthContext"

function ScrollToTop() {
  const { pathname } = useLocation()
  useEffect(() => {
    window.scrollTo(0, 0)
  }, [pathname])
  return null
}

function AppLayout() {
  return (
    <>
      {/* OverlayProvider が overlay を children の後ろに append するので、
          SideMenu は OverlayProvider の外（DOM的に後ろ）に置く。
          これにより z-index なしで「SideMenu > overlay > 通常画面」の重なり順を実現する。 */}
      <OverlayProvider>
        <Header />
        <main className="pt-14">
          <Routes>
            <Route path="/" element={<TopPage />} />
            <Route path="/:slug" element={<MdPage />} />
          </Routes>
        </main>
      </OverlayProvider>
      {/* SideMenu は DOM 最後尾 → overlay より自然に手前に来る */}
      <SideMenu />
    </>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <MenuProvider>
          <ScrollToTop />
          <AppLayout />
        </MenuProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}
