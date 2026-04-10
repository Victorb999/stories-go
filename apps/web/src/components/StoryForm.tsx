import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import type { Story, StoryInput } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useEffect } from 'react'

const schema = z.object({
    title: z.string().min(3, "Título precisa ter no mínimo 3 letras"),
    author: z.string().min(2, "Autor obrigatório"),
    cover_image: z.string().url("Precisa ser uma URL válida"),
    content: z.string().min(10, "Conteúdo muito curto"),
    size: z.enum(['small', 'large']),
    ai_generated: z.boolean(),
})

interface StoryFormProps {
    initialData?: Story;
    onSubmit: (data: StoryInput) => Promise<void>;
    onCancel: () => void;
    isSubmitting: boolean;
}

export function StoryForm({ initialData, onSubmit, onCancel, isSubmitting }: StoryFormProps) {
    const { register, handleSubmit, formState: { errors }, reset, setValue, watch } = useForm<StoryInput>({
        resolver: zodResolver(schema),
        defaultValues: initialData || {
            title: '',
            author: '',
            cover_image: '',
            content: '',
            size: 'small',
            ai_generated: false,
        }
    })

    // Set initial data if it changes
    useEffect(() => {
        if (initialData) {
            reset(initialData)
        }
    }, [initialData, reset])

    const size = watch('size')
    const aiGenerated = watch('ai_generated')

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
                <Label htmlFor="title">Título</Label>
                <Input id="title" {...register('title')} placeholder="Ex: O Menino Mágico" />
                {errors.title && <p className="text-xs text-red-500">{errors.title.message}</p>}
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                    <Label htmlFor="author">Autor</Label>
                    <Input id="author" {...register('author')} placeholder="Ex: Vovó Coruja" />
                    {errors.author && <p className="text-xs text-red-500">{errors.author.message}</p>}
                </div>
                <div className="space-y-2">
                    <Label htmlFor="cover_image">URL da Capa</Label>
                    <Input id="cover_image" {...register('cover_image')} placeholder="https://..." />
                    {errors.cover_image && <p className="text-xs text-red-500">{errors.cover_image.message}</p>}
                </div>
            </div>

            <div className="space-y-2">
                <Label htmlFor="content">Conteúdo da História</Label>
                <Textarea
                    id="content"
                    {...register('content')}
                    placeholder="Era uma vez..."
                    className="h-32 resize-none"
                />
                {errors.content && <p className="text-xs text-red-500">{errors.content.message}</p>}
            </div>

            <div className="flex items-center gap-6 pt-2">
                <div className="space-y-2 flex-1">
                    <Label>Tamanho</Label>
                    <Select
                        value={size}
                        onValueChange={(val: 'small' | 'large') => setValue('size', val, { shouldValidate: true })}
                    >
                        <SelectTrigger>
                            <SelectValue placeholder="Selecione o tamanho" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="small">Curta</SelectItem>
                            <SelectItem value="large">Longa</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                <div className="flex items-center space-x-2 flex-1 pt-6">
                    <Switch
                        id="ai_generated"
                        checked={aiGenerated}
                        onCheckedChange={(val) => setValue('ai_generated', val)}
                    />
                    <Label htmlFor="ai_generated" className="font-normal">Gerada por AI</Label>
                </div>
            </div>

            <div className="flex justify-end gap-2 pt-4">
                <Button type="button" variant="outline" onClick={onCancel} disabled={isSubmitting}>
                    Cancelar
                </Button>
                <Button type="submit" disabled={isSubmitting} className="bg-primary text-primary-foreground hover:bg-primary/90">
                    {isSubmitting ? 'Salvando...' : 'Salvar'}
                </Button>
            </div>
        </form>
    )
}
