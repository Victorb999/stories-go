export interface Story {
    id: number;
    title: string;
    cover_image: string;
    author: string;
    content: string;
    ai_generated: boolean;
    size: 'small' | 'large';
    views: number;
    created_at: string;
    updated_at: string;
}

export type StoryInput = Omit<Story, 'id' | 'views' | 'created_at' | 'updated_at'>;

export async function fetchStories(size?: string, ai?: boolean): Promise<Story[]> {
    const url = new URL('/api/v1/stories', window.location.origin);
    if (size) url.searchParams.set('size', size);
    if (ai !== undefined) url.searchParams.set('ai', ai.toString());

    const res = await fetch(url);
    if (!res.ok) throw new Error('Failed to fetch stories');
    return res.json();
}

export async function fetchStory(id: number): Promise<Story> {
    const res = await fetch(`/api/v1/stories/${id}`);
    if (!res.ok) throw new Error('Failed to fetch story');
    return res.json();
}

export async function createStory(input: StoryInput): Promise<Story> {
    const res = await fetch('/api/v1/stories', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
    });
    if (!res.ok) throw new Error('Failed to create story');
    return res.json();
}

export async function updateStory(id: number, input: StoryInput): Promise<Story> {
    const res = await fetch(`/api/v1/stories/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
    });
    if (!res.ok) throw new Error('Failed to update story');
    return res.json();
}

export async function deleteStory(id: number): Promise<void> {
    const res = await fetch(`/api/v1/stories/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error('Failed to delete story');
}

export async function seedStories(): Promise<void> {
    const res = await fetch('/api/v1/stories/seed', { method: 'POST' });
    if (!res.ok) throw new Error('Failed to seed stories');
}
