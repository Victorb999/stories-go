import { Link, useLocation } from 'react-router-dom'
import { Home, BookOpen, Pencil } from 'lucide-react'
import { cn } from '@/lib/utils'

export function BottomNav() {
    const location = useLocation()

    return (
        <div className="md:hidden fixed bottom-0 left-0 right-0 
        glass-panel border-t border-white/20 flex 
        justify-around items-center py-2 px-6 z-50">
            <Link
                to="/"
                className={cn(
                    "flex flex-col items-center gap-1",
                    location.pathname === '/'
                        ? "text-purple-600 font-bold"
                        : "text-muted-foreground"
                )}
            >
                <Home className="w-5 h-5" />
                <span className="text-[10px]">Início</span>
            </Link>

            <Link
                to="/biblioteca"
                className={cn(
                    "flex flex-col items-center gap-1",
                    location.pathname === '/biblioteca'
                        ? "text-purple-600 font-bold"
                        : "text-muted-foreground"
                )}
            >
                <BookOpen className="w-5 h-5" />
                <span className="text-[10px]">Biblioteca</span>
            </Link>

            <Link
                to="/admin"
                className={cn(
                    "flex flex-col items-center gap-1",
                    location.pathname === '/admin'
                        ? "text-purple-600 font-bold"
                        : "text-muted-foreground"
                )}
            >
                <Pencil className="w-5 h-5" />
                <span className="text-[10px]">Criar</span>
            </Link>

        </div>
    )
}
