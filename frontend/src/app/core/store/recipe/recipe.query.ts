import { Injectable } from "@angular/core";
import { QueryEntity } from "@datorama/akita";
import { RecipeState, RecipeStore } from "./recipe.store";
import { Recipe } from "./recipe.model";
import { Observable } from "rxjs";

@Injectable({
    providedIn: 'root'
})
export class RecipeQuery extends QueryEntity<RecipeState> {
    constructor(protected override store: RecipeStore) {
        super(store);
    }

    selectRecipes() {
        return this.selectAll();
    }

    selectLikedRecipes() {
        return this.select(state => state['likedRecipes'] as Recipe[]);
    }

    selectRecipe(id: string): Observable<Recipe | undefined> {
        return this.selectEntity(id);
    }
}
