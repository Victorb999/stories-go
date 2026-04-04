import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { fetchStory } from '@/lib/api'
import type { Story } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ArrowLeft, Eye, Minus, Plus, Moon, Sun } from 'lucide-react'

export function StoryPage() {
    const { id } = useParams()
    const navigate = useNavigate()
    const [story, setStory] = useState<Story | null>(null)
    const [error, setError] = useState<Error | null>(null)

    // Accessibility states
    const [fontSizeLevel, setFontSizeLevel] = useState(0) // 0 to 4
    const [isNegative, setIsNegative] = useState(false)

    useEffect(() => {
        if (!id) return;
        fetchStory(Number(id))
            .then(setStory)
            .catch(setError)
    }, [id])

    if (error) {
        return <div className="text-center mt-20 text-red-500">Erro ao carregar a história.</div>
    }

    if (!story) {
        return <div className="text-center mt-20 animate-pulse text-primary">Carregando página mágica...</div>
    }

    const sizes = [
        'text-lg sm:text-[1.35rem]',
        'text-xl sm:text-2xl',
        'text-2xl sm:text-3xl',
        'text-3xl sm:text-4xl',
        'text-4xl sm:text-5xl',
    ]
    const currentSize = sizes[fontSizeLevel]

    const bgColor = isNegative ? 'bg-[#0f0e15]' : 'bg-[#FFFDF9]'
    const gradientFrom = isNegative ? 'from-[#0f0e15] via-[#0f0e15]/40' : 'from-[#FFFDF9] via-[#FFFDF9]/40'
    const headingColor = isNegative ? 'text-white' : 'text-[#1A1265]'
    const textColor = isNegative ? 'text-[#E0E2EE]' : 'text-[#2C2E3E]'
    const badgeBg = isNegative ? 'bg-white/10 text-white border-white/20' : 'bg-[#1A1265]/5 text-[#1A1265] border-[#1A1265]'
    const pinkAccent = isNegative ? 'text-primary drop-shadow-[0_0_8px_rgba(249,168,187,0.3)]' : 'text-primary'

    return (
        <div className="max-w-4xl mx-auto animate-in fade-in slide-in-from-bottom-8 duration-700 ease-out mb-16 px-2 sm:px-0">
            <Button
                variant="ghost"
                onClick={() => navigate('/')}
                className="mb-4 text-primary hover:text-white hover:bg-white/10"
            >
                <ArrowLeft className="w-4 h-4 mr-2" />
                Voltar
            </Button>

            <article className={`${bgColor} rounded-3xl sm:rounded-[2.5rem] overflow-hidden shadow-2xl transition-colors duration-500`}>
                <div className="relative">
                    <img
                        src={story.cover_image}
                        alt={story.title}
                        className={`w-full h-[40vh] sm:h-[55vh] object-cover transition-opacity duration-500 ${isNegative ? 'opacity-80' : 'opacity-100'}`}
                    />
                    <div className={`absolute inset-0 bg-gradient-to-t ${gradientFrom} to-transparent h-full w-full`} />
                </div>

                <div className="px-6 sm:px-16 pt-0 pb-20 relative z-10 -mt-20">
                    <div className="flex flex-wrap gap-2 mb-8 justify-center sm:justify-start">
                        {story.ai_generated && (
                            <Badge className="bg-primary text-primary-foreground border-none">
                                Criado por AI
                            </Badge>
                        )}
                        <Badge variant="outline" className={badgeBg}>
                            {story.size === 'small' ? 'História Curta' : 'História Longa'}
                        </Badge>
                        <div className="hidden sm:block flex-1" />
                        <div className={`flex items-center text-sm gap-1.5 px-3 py-1 rounded-full font-bold ${badgeBg}`}>
                            <Eye className="w-4 h-4" />
                            <span>{story.views} lidas</span>
                        </div>
                    </div>

                    <h1 className={`text-4xl sm:text-6xl font-extrabold ${headingColor} mb-4 tracking-tight leading-tight text-center sm:text-left transition-colors duration-500`}>
                        {story.title}
                    </h1>

                    <div className="flex flex-col sm:flex-row items-center justify-between gap-6 mb-14 border-b-2 sm:border-b-0 border-[#F9A8BB]/30 pb-6 sm:pb-0">
                        <p className={`text-xl font-bold ${pinkAccent} drop-shadow-sm transition-colors duration-500`}>
                            Por <span className={headingColor}>{story.author}</span>
                        </p>

                        {/* Accessibility Tools */}
                        <div className={`flex items-center gap-1.5 p-1.5 rounded-full border ${isNegative ? 'bg-white/5 border-white/10' : 'bg-black/5 border-black/10'}`}>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setFontSizeLevel(Math.max(0, fontSizeLevel - 1))}
                                disabled={fontSizeLevel === 0}
                                className={`w-8 h-8 rounded-full ${isNegative ? 'text-white hover:bg-white/20' : 'text-black hover:bg-black/10'}`}
                            >
                                <Minus className="w-3.5 h-3.5" />
                            </Button>
                            <span className={`text-xs font-bold w-6 text-center ${headingColor}`}>{fontSizeLevel + 1}</span>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setFontSizeLevel(Math.min(sizes.length - 1, fontSizeLevel + 1))}
                                disabled={fontSizeLevel === sizes.length - 1}
                                className={`w-8 h-8 rounded-full mr-2 ${isNegative ? 'text-white hover:bg-white/20' : 'text-black hover:bg-black/10'}`}
                            >
                                <Plus className="w-3.5 h-3.5" />
                            </Button>

                            <div className={`w-px h-6 mx-1 ${isNegative ? 'bg-white/20' : 'bg-black/20'}`} />

                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setIsNegative(!isNegative)}
                                className={`w-9 h-9 rounded-full ${isNegative ? 'text-yellow-300 hover:bg-white/20' : 'text-indigo-900 hover:bg-black/10'}`}
                            >
                                {isNegative ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
                            </Button>
                        </div>
                    </div>

                    <div className="mx-auto max-w-2xl space-y-8">
                        {story.content.split('\n\n').map((paragraph, i) => {
                            if (!paragraph.trim()) return null;
                            const isFirst = i === 0;
                            return (
                                <p
                                    key={i}
                                    className={`font-serif leading-relaxed sm:leading-loose text-justify ${currentSize} ${textColor} tracking-wide transition-all duration-500
                                ${isFirst ? `first-letter:text-[3em] sm:first-letter:text-[3.5em] first-letter:font-bold first-letter:${pinkAccent} first-letter:mr-2 first-letter:float-left` : ''}`}
                                >
                                    {paragraph}
                                </p>
                            )
                        })}
                    </div>
                </div>
            </article>
        </div>
    )
}
