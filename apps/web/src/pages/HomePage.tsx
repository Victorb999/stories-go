import { useAtom, useAtomValue } from 'jotai'
import { useNavigate } from 'react-router-dom'
import { storiesAtom, sizeFilterAtom, aiFilterAtom } from '@/store/atoms'
import { StoryCard } from '@/components/StoryCard'
import { Badge } from '@/components/ui/badge'
import { Rocket, ChevronRight } from 'lucide-react'

export function HomePage() {
    const stories = useAtomValue(storiesAtom)
    const [size, setSize] = useAtom(sizeFilterAtom)
    const [ai, setAi] = useAtom(aiFilterAtom)
    const navigate = useNavigate()

    // Pick featured stories for the bento grid
    const featured = stories.slice(0, 4)
    const mainStory = featured[0]
    const secondStory = featured[1]
    const smallStories = featured.slice(2, 4)

    return (
        <div className="space-y-16 animate-in fade-in slide-in-from-bottom-4 duration-500 ease-out">
            {/* ── Hero Section ── */}
            <section className="relative min-h-[500px] flex items-center rounded-2xl overflow-hidden mt-4">
                <div className="absolute inset-0 z-0">
                    <img
                        src="/hero.png"
                        alt="Cena mágica de leitura"
                        className="w-full h-full object-cover opacity-80"
                    />
                    <div className="absolute inset-0 bg-gradient-to-r from-secondary-container via-secondary-container/60 to-transparent" />
                </div>

                <div className="relative z-10 max-w-2xl px-8 md:px-12 space-y-6">
                    <h1 className="text-4xl md:text-7xl font-extrabold text-on-secondary-container leading-tight tracking-tight">
                        Onde Cada Página{' '}
                        <span className="text-primary">Respira Magia</span>
                    </h1>
                    <p className="text-lg md:text-xl text-on-secondary-container/80 font-medium leading-relaxed">
                        Contos de fadas personalizados para seus pequenos exploradores. Solte a
                        imaginação com histórias criadas por IA onde eles são os heróis.
                    </p>
                    <div className="flex flex-wrap gap-4 pt-4">
                        <button
                            onClick={() => navigate('/admin')}
                            className="bg-primary text-primary-foreground px-8 py-4 rounded-full text-lg font-bold shadow-lg shadow-primary/20 hover:scale-105 active:scale-95 transition-all flex items-center gap-2"
                        >
                            Criar Sua História <Rocket className="w-5 h-5" />
                        </button>
                        <button
                            onClick={() => navigate('/biblioteca')}
                            className="bg-white text-secondary px-8 py-4 rounded-full text-lg font-bold hover:bg-secondary-container transition-all flex items-center gap-2"
                        >
                            Explorar Biblioteca
                        </button>
                    </div>
                </div>
            </section>

            {/* ── Histórias Populares Bento Grid ── */}
            <section id="popular-tales" className="space-y-8">
                <div className="flex flex-col md:flex-row justify-between items-end gap-4 px-2">
                    <div className="space-y-2">
                        <h2 className="text-3xl md:text-4xl font-bold text-foreground tracking-tight">
                            Histórias Populares
                        </h2>
                        <p className="text-muted-foreground font-medium">
                            As mais mágicas, votadas pela nossa comunidade de sonhadores.
                        </p>
                    </div>
                    <button
                        onClick={() => navigate('/biblioteca')}
                        className="text-primary font-bold flex items-center gap-1 hover:underline"
                    >
                        Ver todas as histórias <ChevronRight className="w-4 h-4" />
                    </button>
                </div>

                {/* Filter Chips */}
                <div className="flex items-center gap-2 overflow-x-auto pb-1 scrollbar-hide">
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

                {stories.length === 0 ? (
                    <div className="py-20 text-center flex flex-col items-center">
                        <div className="text-4xl mb-4 drop-shadow-md">📚</div>
                        <h3 className="text-lg font-medium text-foreground">Nenhuma história encontrada</h3>
                        <p className="text-muted-foreground mt-1 text-sm">
                            Vá para o painel Admin e use o botão "Seed" para gerar histórias de exemplo.
                        </p>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                        {/* Large Feature Card */}
                        {mainStory && (
                            <div
                                onClick={() => navigate(`/stories/${mainStory.id}`)}
                                className="md:col-span-2 md:row-span-2 group relative overflow-hidden rounded-2xl bg-card border-4 border-primary-container p-2 transition-transform hover:-rotate-1 cursor-pointer"
                            >
                                <div className="relative h-64 md:h-[480px] rounded-xl overflow-hidden">
                                    <img
                                        src={mainStory.cover_image}
                                        alt={mainStory.title}
                                        className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                    />
                                    <div className="absolute inset-0 bg-gradient-to-t from-primary/90 via-transparent to-transparent" />
                                    <div className="absolute bottom-0 p-8 text-primary-foreground">
                                        <span className="bg-white/20 backdrop-blur-md px-3 py-1 rounded-full text-xs font-bold uppercase tracking-widest mb-4 inline-block">
                                            Mais Lida
                                        </span>
                                        <h3 className="text-2xl md:text-3xl font-bold mb-2">{mainStory.title}</h3>
                                        <p className="text-primary-foreground/80 line-clamp-2 text-sm">
                                            Por {mainStory.author}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Horizontal Card */}
                        {secondStory && (
                            <div
                                onClick={() => navigate(`/stories/${secondStory.id}`)}
                                className="md:col-span-2 group relative overflow-hidden rounded-2xl bg-card border-4 border-tertiary-container p-2 transition-transform hover:rotate-1 cursor-pointer"
                            >
                                <div className="relative h-64 rounded-xl overflow-hidden flex flex-col md:flex-row">
                                    <div className="w-full md:w-1/2 h-full">
                                        <img
                                            src={secondStory.cover_image}
                                            alt={secondStory.title}
                                            className="w-full h-full object-cover"
                                        />
                                    </div>
                                    <div className="w-full md:w-1/2 p-6 flex flex-col justify-center bg-tertiary-container/10">
                                        <h3 className="text-xl font-bold text-on-tertiary-container mb-2">
                                            {secondStory.title}
                                        </h3>
                                        <p className="text-on-tertiary-container/70 text-sm mb-4">
                                            Por {secondStory.author}
                                        </p>
                                        <span className="bg-tertiary text-on-tertiary w-fit px-4 py-2 rounded-full text-sm font-bold">
                                            Ler Agora
                                        </span>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Small Cards */}
                        {smallStories.map((story, i) => (
                            <div
                                key={story.id}
                                onClick={() => navigate(`/stories/${story.id}`)}
                                className={`group relative overflow-hidden rounded-2xl bg-card border-4 ${i === 0 ? 'border-outline-soft' : 'border-secondary-container'
                                    } p-2 transition-transform cursor-pointer ${i === 0 ? 'hover:-translate-y-2' : 'hover:translate-y-2'
                                    }`}
                            >
                                <div className="relative h-48 rounded-xl overflow-hidden mb-4">
                                    <img
                                        src={story.cover_image}
                                        alt={story.title}
                                        className="w-full h-full object-cover"
                                    />
                                </div>
                                <div className="px-2 pb-2">
                                    <h3 className="font-bold text-foreground">{story.title}</h3>
                                    <p className="text-xs text-muted-foreground">
                                        Por {story.author}
                                    </p>
                                </div>
                            </div>
                        ))}

                        {/* Remaining stories as standard cards */}
                        {stories.slice(4).map(story => (
                            <div key={story.id} className="snap-center shrink-0 h-full">
                                <StoryCard
                                    story={story}
                                    onClick={() => navigate(`/stories/${story.id}`)}
                                />
                            </div>
                        ))}
                    </div>
                )}
            </section>
        </div>
    )
}
