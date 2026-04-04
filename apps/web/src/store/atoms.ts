import { atom } from 'jotai'
import { fetchStories } from '../lib/api'
import type { Story } from '../lib/api'

// Atom for the refresh trigger
export const refreshStoriesTriggerAtom = atom(0)

// Primitive atoms for filters
export const sizeFilterAtom = atom<string | undefined>(undefined)
export const aiFilterAtom = atom<boolean | undefined>(undefined)

// Async atom to fetch stories based on filters and refresh trigger
export const storiesAtom = atom<Promise<Story[]>>(async (get) => {
    // Read trigger so it re-fetches when incremented
    get(refreshStoriesTriggerAtom)

    const size = get(sizeFilterAtom)
    const ai = get(aiFilterAtom)
    return fetchStories(size, ai)
})
