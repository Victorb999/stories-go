export function Footer() {
    return (
        <footer className="bg-muted flex flex-col items-center gap-4 py-10 px-8 w-full rounded-t-[3rem] mt-12 transition-all duration-300">
            <div className="flex flex-col items-center gap-2">
                <span className="font-bold text-primary text-xl tracking-tight">BabyStories</span>
                <p className="text-muted-foreground text-sm text-center max-w-xs">
                    © 2024 BabyStories. Magia em cada página.
                </p>
            </div>
        </footer>
    )
}
