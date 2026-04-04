import { Link, useLocation } from 'react-router-dom'
import { BookOpen, Settings } from 'lucide-react'
import { cn } from '@/lib/utils'

export function BottomNav() {
    const location = useLocation()

    return (
        <div className="fixed bottom-0 mt-20 left-0 right-0 z-50 bg-[#140D4F]/90 backdrop-blur-lg border-t border-border sm:hidden pb-safe">
            <div className="flex items-center justify-around p-2">
                <Link
                    to="/"
                    className={cn(
                        "flex flex-col items-center gap-1 p-2 rounded-xl transition-colors min-w-[64px]",
                        location.pathname === '/' || location.pathname.startsWith('/stories')
                            ? "text-primary"
                            : "text-muted-foreground hover:text-white"
                    )}
                >
                    <BookOpen className="w-5 h-5" />
                    <span className="text-[10px] font-medium">Histórias</span>
                </Link>
                <Link
                    to="/admin"
                    className={cn(
                        "flex flex-col items-center gap-1 p-2 rounded-xl transition-colors min-w-[64px]",
                        location.pathname === '/admin'
                            ? "text-primary"
                            : "text-muted-foreground hover:text-white"
                    )}
                >
                    <Settings className="w-5 h-5" />
                    <span className="text-[10px] font-medium">Admin</span>
                </Link>
            </div>
        </div>
    )
}
