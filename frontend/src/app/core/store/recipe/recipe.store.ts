import { Injectable } from "@angular/core";
import { EntityState, EntityStore, StoreConfig } from '@datorama/akita';
import { Recipe } from "./recipe.model";

// Step 1: Extend the state with your custom field
export interface RecipeState extends EntityState<Recipe> {
    likedRecipeIds: Set<string>;  // Add this!
}

@Injectable({
    providedIn: 'root'
})
@StoreConfig({ name: 'recipe' })
export class RecipeStore extends EntityStore<RecipeState> {
    constructor() {
        super({ likedRecipeIds: new Set<string>() }); // Initialize it!
    }

    // Add method
    addToLiked(recipeId: string) {
        this.update(state => ({
            likedRecipeIds: new Set(state.likedRecipeIds).add(recipeId)
        }));
    }

    removeFromLiked(recipeId: string) {
        this.update(state => {
            const newSet = new Set(state.likedRecipeIds);
            newSet.delete(recipeId);
            return { likedRecipeIds: newSet };
        });
    }

    // Optional: set full liked recipes
    setLikedRecipes(recipes: Recipe[]) {
        this.update(state => ({
            likedRecipeIds: new Set(recipes.map(r => r.id))
        }));
        this.upsertMany(recipes);
    }

    // Update a single recipe's like state
    updateLike(id: string, changes: { likes_count: number; is_liked: boolean }) {
        this.update(id, changes);
    }

}