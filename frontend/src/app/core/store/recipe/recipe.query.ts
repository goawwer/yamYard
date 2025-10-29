import { Injectable } from "@angular/core";
import { QueryEntity } from "@datorama/akita";
import { RecipeState, RecipeStore } from "./recipe.store";

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
}
