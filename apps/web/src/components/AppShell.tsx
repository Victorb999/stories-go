import { BottomNav } from './BottomNav'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { BookOpen } from 'lucide-react'
import { Suspense } from 'react'

export function AppShell() {
    const location = useLocation()

    return (
        <div className="min-h-screen bg-background flex flex-col font-sans pb-16 sm:pb-0 text-foreground">
            <header className="sticky top-0 z-50 bg-[#140D4F]/90 backdrop-blur-md border-b border-border">
                <div className="container mx-auto px-4 h-14 flex items-center justify-between">
                    <Link to="/" className="flex items-center gap-2 text-primary hover:opacity-80 transition-opacity">
                        <div className="bg-primary p-1.5 rounded-lg text-primary-foreground">
                            <BookOpen className="w-5 h-5" />
                        </div>
                        <span className="font-bold text-lg tracking-tight">BabyStories</span>
                    </Link>

                    <nav className="hidden sm:flex items-center gap-6">
                        <Link
                            to="/"
                            className={`text-sm font-bold transition-colors hover:text-primary ${location.pathname === '/' ? 'text-primary' : 'text-primary/60'}`}
                        >
                            Histórias
                        </Link>
                        <Link
                            to="/admin"
                            className={`text-sm font-bold transition-colors hover:text-primary ${location.pathname === '/admin' ? 'text-primary' : 'text-primary/60'}`}
                        >
                            Admin
                        </Link>
                    </nav>
                </div>
            </header>

            <main className="flex-1 container mx-auto px-4 py-8">
                <Suspense fallback={<div className="flex h-[50vh] items-center justify-center text-primary"><div className="animate-pulse">Carregando...</div></div>}>
                    <Outlet />
                </Suspense>
            </main>

            <BottomNav />
        </div>
    )
}
