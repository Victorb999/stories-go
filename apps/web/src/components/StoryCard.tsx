import type { Story } from '@/lib/api'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Eye } from 'lucide-react'

export function StoryCard({ story, onClick }: { story: Story; onClick: () => void }) {
    return (
        <Card
            onClick={onClick}
            className="group relative overflow-hidden cursor-pointer transition-all hover:scale-[1.02] hover:shadow-primary/20 hover:shadow-xl border-2 border-border aspect-[3/4] rounded-2xl"
        >
            <img
                src={story.cover_image}
                alt={story.title}
                className="absolute inset-0 w-full h-full object-cover transition-transform group-hover:scale-110 duration-700"
                loading="lazy"
            />
            {/* Gradient overlay for readability */}
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent opacity-90 transition-opacity group-hover:opacity-100" />

            <div className="absolute inset-0 flex flex-col justify-end p-5">
                <div className="flex flex-wrap gap-2 mb-3">
                    {story.ai_generated && (
                        <Badge className="bg-primary text-primary-foreground hover:bg-primary/90 text-[10px] uppercase font-bold tracking-wider px-2 py-0.5 border-none">
                            AI
                        </Badge>
                    )}
                    <Badge variant="outline" className="text-[10px] uppercase font-bold tracking-wider px-2 py-0.5 border-white/40 text-white backdrop-blur-sm">
                        {story.size === 'small' ? 'Curta' : 'Longa'}
                    </Badge>
                </div>
                <h3 className="text-xl font-bold line-clamp-2 text-white mb-2 leading-tight">{story.title}</h3>
                <div className="flex items-center justify-between mt-2">
                    <p className="text-sm font-bold text-primary-container line-clamp-1">Por {story.author}</p>
                    <span className="flex items-center gap-1.5 text-xs text-white/80 font-bold whitespace-nowrap"><Eye className="w-3.5 h-3.5" /> {story.views}</span>
                </div>
            </div>
        </Card>
    )
}
