import { Injectable } from "@angular/core";
import { EntityState, EntityStore, StoreConfig } from '@datorama/akita';
import { Recipe } from "./recipe.model";

export interface RecipeState extends EntityState<Recipe> { }

@Injectable({
    providedIn: 'root'
})
@StoreConfig({ name: 'recipe' })
export class RecipeStore extends EntityStore<RecipeState> {
    constructor() {
        super();
    }
}