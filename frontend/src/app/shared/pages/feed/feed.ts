import { Component, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { Recipe } from '../../../core/store/recipe/recipe.model';
import { RecipeService } from '../../../core/store/recipe/recipe.service';
import { RecipeStore } from '../../../core/store/recipe/recipe.store';
import { FormControl, FormGroup } from '@angular/forms';
import { RecipeQuery } from '../../../core/store/recipe/recipe.query';
import { UserService } from '../../../core/store/user/user.service';
import { Sort } from '@angular/material/sort';
import { UserQuery } from '../../../core/store/user/user.query';
import { FEEDCOMPONENTS } from './feed.imports';

@Component({
    selector: 'app-feed',
    imports: [FEEDCOMPONENTS],
    templateUrl: './feed.html',
    styleUrl: './feed.scss'
})
export class FeedComponent {
    isCreating = signal(false);
    currentSort = signal<Sort | null>(null);
    filterForm: FormGroup = new FormGroup({
        username: new FormControl(''),
    });
    protected recipes$!: Observable<Recipe[]>;
    recipeForm = new FormGroup({
        title: new FormControl(''),
        description: new FormControl(''),
        ingredients: new FormControl(''),
        cookingTime: new FormControl<number | null>(null),
        difficulty: new FormControl(''),
    });

    constructor(
        private service: RecipeService,
        private query: RecipeQuery,
        private userService: UserService,
        private userQuery: UserQuery
    ) {
    }

    ngOnInit(): void {
        this.loadData();
        this.recipes$ = this.query.selectRecipes();
    }

    private loadData() {
        const sort = this.currentSort();
        this.service.getAllRecipes(sort ?? undefined, undefined).subscribe();
    }

    sortData(sort: Sort) {
        if (!sort.direction) this.currentSort.set(null);
        else this.currentSort.set(sort);

        this.loadData();
    }

    toggleCreateForm() {
        this.isCreating.update((v) => !v);
    }

    createRecipe() {
        const activeUser = this.userQuery.getActive();
        if (!activeUser) return;

        const now = new Date().toISOString();

        const recipe: Recipe = {
            id: crypto.randomUUID(),
            authorId: activeUser.id,
            authorUsername: activeUser.username,
            title: this.recipeForm.value.title!,
            description: this.recipeForm.value.description!,
            ingredients: this.recipeForm.value.ingredients?.split(',').map(s => s.trim()) || [],
            cookingTime: this.recipeForm.value.cookingTime || 0,
            difficulty: this.recipeForm.value.difficulty || 'Medium',
            createdAt: now,
            updatedAt: now,
        };

        this.service.postRecipe(recipe).subscribe(() => {
            this.toggleCreateForm();
            this.service.getAllRecipes().subscribe(); // reload feed
            this.recipeForm.reset();
        });
    }
}
