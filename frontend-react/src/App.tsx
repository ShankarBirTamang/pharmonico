import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Header, Footer } from './components/layout'
import { Dashboard, Enroll, Home } from './pages'

// Layout wrapper component for routes that need header/footer
function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        {children}
      </main>
      <Footer />
    </>
  )
}

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-teal-50 dark:from-slate-950 dark:via-slate-900 dark:to-slate-950">
        <Routes>
          {/* Enrollment route - no header/footer */}
          <Route path="/enroll/:token" element={<Enroll />} />
          
          {/* Routes with layout */}
          <Route
            path="/"
            element={
              <Layout>
                <Home />
              </Layout>
            }
          />
          <Route
            path="/ops/dashboard"
            element={
              <Layout>
                <Dashboard />
              </Layout>
            }
          />
          
          {/* Catch-all redirect */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default App
