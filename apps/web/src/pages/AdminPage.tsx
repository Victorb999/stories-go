import { useSetAtom, useAtomValue } from 'jotai'
import { useState, useTransition } from 'react'
import { storiesAtom, refreshStoriesTriggerAtom } from '@/store/atoms'
import { createStory, updateStory, deleteStory, seedStories } from '@/lib/api'
import type { Story, StoryInput } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { StoryForm } from '@/components/StoryForm'
import { Plus, Trash2, Edit2, Database } from 'lucide-react'

export function AdminPage() {
    const stories = useAtomValue(storiesAtom)
    const refresh = useSetAtom(refreshStoriesTriggerAtom)
    const [isPending, startTransition] = useTransition()

    const [isFormOpen, setIsFormOpen] = useState(false)
    const [editingStory, setEditingStory] = useState<Story | undefined>(undefined)

    const handleCreate = async (data: StoryInput) => {
        await createStory(data)
        setIsFormOpen(false)
        startTransition(() => {
            refresh(prev => prev + 1)
        })
    }

    const handleUpdate = async (data: StoryInput) => {
        if (!editingStory) return
        await updateStory(editingStory.id, data)
        setIsFormOpen(false)
        startTransition(() => {
            refresh(prev => prev + 1)
        })
    }

    const handleDelete = async (id: number) => {
        if (!confirm('Tem certeza?')) return
        await deleteStory(id)
        startTransition(() => {
            refresh(prev => prev + 1)
        })
    }

    const handleSeed = async () => {
        await seedStories()
        startTransition(() => {
            refresh(prev => prev + 1)
        })
    }

    return (
        <div className="space-y-8 animate-in fade-in duration-500">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white">Gerenciamento</h1>
                    <p className="text-primary/80 mt-1 text-sm">Adicione, edite ou remova histórias.</p>
                </div>

                <div className="flex gap-2">
                    {stories.length === 0 && (
                        <Button variant="outline" onClick={handleSeed} disabled={isPending} className="border-[#F9A8BB] text-[#F9A8BB] hover:bg-[#F9A8BB]/10">
                            <Database className="w-4 h-4 mr-2" />
                            Preencher com Exemplo
                        </Button>
                    )}

                    <Dialog open={isFormOpen} onOpenChange={(open) => {
                        setIsFormOpen(open)
                        if (!open) setEditingStory(undefined)
                    }}>
                        <DialogTrigger asChild>
                            <Button className="bg-primary text-primary-foreground hover:bg-primary/90">
                                <Plus className="w-4 h-4 mr-2" />
                                Nova História
                            </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-[600px] bg-[#140D4F] text-white border-border">
                            <DialogHeader>
                                <DialogTitle>{editingStory ? 'Editar História' : 'Nova História'}</DialogTitle>
                            </DialogHeader>
                            <StoryForm
                                initialData={editingStory}
                                onSubmit={editingStory ? handleUpdate : handleCreate}
                                onCancel={() => setIsFormOpen(false)}
                                isSubmitting={isPending}
                            />
                        </DialogContent>
                    </Dialog>
                </div>
            </div>

            <div className="bg-[#140D4F] rounded-2xl shadow-sm border border-border overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm">
                        <thead className="bg-[#201773] text-primary/90 font-medium">
                            <tr>
                                <th className="px-6 py-4">ID</th>
                                <th className="px-6 py-4">Título</th>
                                <th className="px-6 py-4">Autor</th>
                                <th className="px-6 py-4">Tamanho</th>
                                <th className="px-6 py-4">AI</th>
                                <th className="px-6 py-4 text-right">Ações</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-border">
                            {stories.length === 0 ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-8 text-center text-muted-foreground">
                                        Nenhuma história. Clique em "Nova História" ou adicione exemplos.
                                    </td>
                                </tr>
                            ) : (
                                stories.map(story => (
                                    <tr key={story.id} className="hover:bg-[#201773]/50 transition-colors">
                                        <td className="px-6 py-4 text-white/50">#{story.id}</td>
                                        <td className="px-6 py-4 font-medium text-white">{story.title}</td>
                                        <td className="px-6 py-4 text-white/80">{story.author}</td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2 py-1 rounded-md text-xs font-semibold ${story.size === 'small' ? 'bg-primary/20 text-primary' : 'bg-green-500/20 text-green-400'}`}>
                                                {story.size === 'small' ? 'Curta' : 'Longa'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-white/80">
                                            {story.ai_generated ? '✨ Sim' : 'Não'}
                                        </td>
                                        <td className="px-6 py-4 text-right space-x-2 whitespace-nowrap">
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => {
                                                    setEditingStory(story)
                                                    setIsFormOpen(true)
                                                }}
                                                className="hover:bg-primary/20 hover:text-primary"
                                            >
                                                <Edit2 className="w-4 h-4 text-primary" />
                                            </Button>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => handleDelete(story.id)}
                                                disabled={isPending}
                                                className="hover:bg-red-500/20 hover:text-red-400 text-red-500"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </Button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    )
}
