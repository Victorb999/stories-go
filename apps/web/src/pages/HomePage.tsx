import { useAtom, useAtomValue } from 'jotai'
import { useNavigate } from 'react-router-dom'
import { storiesAtom, sizeFilterAtom, aiFilterAtom } from '@/store/atoms'
import { StoryCard } from '@/components/StoryCard'
import { Badge } from '@/components/ui/badge'

export function HomePage() {
    const stories = useAtomValue(storiesAtom)
    const [size, setSize] = useAtom(sizeFilterAtom)
    const [ai, setAi] = useAtom(aiFilterAtom)
    const navigate = useNavigate()

    return (
        <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500 ease-out">
            <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white drop-shadow-sm">Histórias Infantis</h1>
                    <p className="text-primary/90 mt-1 text-sm font-medium">Explore contos mágicos para ninar.</p>
                </div>

                <div className="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0">
                    <Badge
                        variant={size === undefined && ai === undefined ? "default" : "secondary"}
                        className="cursor-pointer whitespace-nowrap"
                        onClick={() => { setSize(undefined); setAi(undefined) }}
                    >
                        Todas
                    </Badge>
                    <Badge
                        variant={size === 'small' ? "default" : "secondary"}
                        className="cursor-pointer whitespace-nowrap"
                        onClick={() => setSize(size === 'small' ? undefined : 'small')}
                    >
                        Curtas
                    </Badge>
                    <Badge
                        variant={size === 'large' ? "default" : "secondary"}
                        className="cursor-pointer whitespace-nowrap"
                        onClick={() => setSize(size === 'large' ? undefined : 'large')}
                    >
                        Longas
                    </Badge>
                    <Badge
                        variant={ai === true ? "default" : "secondary"}
                        className="cursor-pointer whitespace-nowrap"
                        onClick={() => setAi(ai === true ? undefined : true)}
                    >
                        Com AI
                    </Badge>
                </div>
            </div>

            {stories.length === 0 ? (
                <div className="py-20 text-center flex flex-col items-center">
                    <div className="text-4xl mb-4 drop-shadow-md">📚</div>
                    <h3 className="text-lg font-medium text-white">Nenhuma história encontrada</h3>
                    <p className="text-primary/70 mt-1 text-sm">Vá para o painel Admin e use o botão "Seed" para gerar histórias de exemplo.</p>
                </div>
            ) : (
                <div className="flex sm:grid sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 sm:gap-6 overflow-x-auto sm:overflow-visible snap-x snap-mandatory pb-6 pt-2 -mx-4 px-4 sm:mx-0 sm:px-0 scrollbar-hide">
                    {stories.map(story => (
                        <div key={story.id} className="snap-center shrink-0 w-[70vw] sm:w-auto h-full">
                            <StoryCard
                                story={story}
                                onClick={() => navigate(`/stories/${story.id}`)}
                            />
                        </div>
                    ))}
                </div>
            )}
        </div>
    )
}
