import { BottomNav } from './BottomNav'
import { Footer } from './Footer'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { Sparkles } from 'lucide-react'
import { Suspense } from 'react'

export function AppShell() {
    const location = useLocation()

    return (
        <div className="min-h-screen flex flex-col font-sans text-foreground">
            {/* Top App Bar */}
            <header className="bg-background/90 backdrop-blur-md flex justify-between items-center px-6 py-4 w-full fixed top-0 z-50 transition-all duration-300 border-b border-border/40">
                <Link to="/" className="flex items-center gap-2">
                    <span className="text-2xl font-extrabold bg-gradient-to-r from-purple-500 to-pink-400 bg-clip-text text-transparent tracking-tight">
                        BabyStories
                    </span>
                </Link>

                <nav className="hidden md:flex items-center gap-8">
                    <Link
                        to="/"
                        className={`font-bold transition-opacity hover:opacity-80 ${location.pathname === '/'
                                ? 'text-purple-600 border-b-2 border-purple-500'
                                : 'text-muted-foreground'
                            }`}
                    >
                        Início
                    </Link>
                    <Link
                        to="/biblioteca"
                        className={`font-bold transition-opacity hover:opacity-80 ${location.pathname === '/biblioteca'
                                ? 'text-purple-600 border-b-2 border-purple-500'
                                : 'text-muted-foreground'
                            }`}
                    >
                        Biblioteca
                    </Link>
                    <Link
                        to="/admin"
                        className={`font-bold transition-opacity hover:opacity-80 ${location.pathname === '/admin'
                                ? 'text-purple-600 border-b-2 border-purple-500'
                                : 'text-muted-foreground'
                            }`}
                    >
                        Criar
                    </Link>
                </nav>

                <div className="flex items-center gap-4">
                    <button className="text-purple-500 hover:opacity-80 transition-opacity active:scale-95 duration-200 ease-out">
                        <Sparkles className="w-5 h-5" />
                    </button>
                    <Link
                        to="/biblioteca"
                        className="bg-primary text-primary-foreground px-6 py-2 rounded-full font-bold hover:opacity-90 active:scale-95 transition-all text-sm"
                    >
                        Começar a Ler
                    </Link>
                </div>
            </header>

            {/* Main Content */}
            <main className="flex-1 pt-20 pb-12 px-6 max-w-7xl mx-auto w-full space-y-16">
                <Suspense
                    fallback={
                        <div className="flex h-[50vh] items-center justify-center text-primary">
                            <div className="animate-pulse text-lg">Carregando...</div>
                        </div>
                    }
                >
                    <Outlet />
                </Suspense>
            </main>

            <Footer />
            <BottomNav />
        </div>
    )
}
