import { Routes, Route } from 'react-router-dom'
import { AppShell } from '@/components/AppShell'
import { HomePage } from '@/pages/HomePage'
import { StoryPage } from '@/pages/StoryPage'
import { AdminPage } from '@/pages/AdminPage'

function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/stories/:id" element={<StoryPage />} />
        <Route path="/admin" element={<AdminPage />} />
      </Route>
    </Routes>
  )
}

export default App
