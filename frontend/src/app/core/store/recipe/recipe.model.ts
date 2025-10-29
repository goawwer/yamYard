export interface Recipe {
    id: string;
    authorId: string;
    authorUsername: string;
    title: string;
    description: string;
    ingredients: string[];
    cookingTime?: number;
    difficulty?: string;
    imageURL?: string;
    createdAt: string;
    updatedAt: string;
}