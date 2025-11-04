export interface Recipe {
    id: string;
    author_id: string;
    author_username: string;
    title: string;
    description: string;
    ingredients: string;
    cooking_time?: number;
    difficulty?: string;
    image_url?: string;
    created_at: string;
    updated_at: string;
}

export interface isUpdatingRecipe {
    title: string | null | undefined;
    description: string | null | undefined;
    ingredients: string | null | undefined;
    cooking_time: number | null | undefined;
    difficulty: string | null | undefined;
    image_url: string | null | undefined;
}