import { Component, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';
import { Recipe } from '../../../core/store/recipe/recipe.model';
import { RecipeService } from '../../../core/store/recipe/recipe.service';
import { RecipeQuery } from '../../../core/store/recipe/recipe.query';
import { UserQuery } from '../../../core/store/user/user.query';
import { FormControl, FormGroup } from '@angular/forms';
import { Sort } from '@angular/material/sort';
import { UserService } from '../../../core/store/user/user.service';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { FEEDCOMPONENTS } from './feed.imports';
import { User } from '../../../core/store/user/user.model';
import { RecipeStore } from '../../../core/store/recipe/recipe.store';

@Component({
    selector: 'app-feed',
    imports: [FEEDCOMPONENTS],
    templateUrl: './feed.html',
    styleUrl: './feed.scss'
})
export class FeedComponent {
    // ── YOUR EXISTING LOGIC ─────────────────────────────────────
    isCreating = signal(false);
    currentSort = signal<Sort | null>(null);
    currentFilter = signal<{ column: string; value: string } | undefined>(undefined);
    currentAvatarUrl = signal<string | null>(null);
    user$!: Observable<User | undefined>;
    filterForm: FormGroup = new FormGroup({
        username: new FormControl(''),
    });

    protected recipes$!: Observable<Recipe[]>;
    recipeForm = new FormGroup({
        title: new FormControl(''),
        description: new FormControl(''),
        ingredients: new FormControl(''),
        cooking_time: new FormControl<number | null>(null),
        difficulty: new FormControl(''),
        image_url: new FormControl(null),
    });

    // ── SEARCH CONTROL ──────────────────────────────────────────
    searchCtrl = new FormControl('', { nonNullable: true });

    // ── MODERN FEATURES ─────────────────────────────────────────
    previewUrl = signal<SafeUrl | null>(null);
    selectedFile?: File;
    isExpanded = new Set<string>();


    constructor(
        private service: RecipeService,
        private query: RecipeQuery,
        private store: RecipeStore,
        private userService: UserService,
        private userQuery: UserQuery,
        private sanitizer: DomSanitizer,
    ) { }

    ngOnInit(): void {
        // ── SEARCH SUBSCRIPTION ────────────────────────────────────
        this.searchCtrl.valueChanges.pipe(
            debounceTime(300),
            distinctUntilChanged()
        ).subscribe(term => {
            this.currentFilter.set(term.trim() ? { column: 'title', value: term.trim() } : undefined);
            this.loadData();
        });

        this.loadData();
        this.recipes$ = this.query.selectRecipes();

        this.user$ = this.userQuery.selectAuthenticated();

        this.user$.subscribe(user => {
            if (user?.image_url) {
                this.currentAvatarUrl.set(user.image_url);
            }
        });
    }

    // ── YOUR UPDATED LOAD DATA (with filter support) ────────────
    private loadData() {
        this.service.getAllRecipes(this.currentSort() ?? undefined, this.currentFilter()).subscribe();
    }

    // ── YOUR SORT LOGIC ─────────────────────────────────────────
    sortData(sort: Sort) {
        if (!sort.direction) this.currentSort.set(null);
        else this.currentSort.set(sort);
        this.loadData();
    }

    onSortChange(sortKey: string) {
        let sort: Sort | null = null;
        switch (sortKey) {
            case 'newest':
                sort = { active: 'createdAt', direction: 'desc' };
                break;
            case 'oldest':
                sort = { active: 'createdAt', direction: 'asc' };
                break;
        }
        this.currentSort.set(sort);
        this.loadData();
    }

    // ── YOUR CREATE LOGIC ───────────────────────────────────────
    toggleCreateForm() {
        this.isCreating.update((v) => !v);
        if (!this.isCreating()) {
            this.recipeForm.reset();
            this.previewUrl.set(null);
            this.selectedFile = undefined;
        }
    }

    createRecipe() {
        const activeUser = this.userQuery.getActive();
        if (!activeUser) return;

        const now = new Date().toISOString();

        const recipe: Recipe = {
            id: crypto.randomUUID(),
            author_id: activeUser.id,
            author_username: activeUser.username,
            title: this.recipeForm.value.title!,
            description: this.recipeForm.value.description!,
            ingredients: this.recipeForm.value.ingredients!,
            cooking_time: this.recipeForm.value.cooking_time || 0,
            difficulty: this.recipeForm.value.difficulty?.toLocaleLowerCase() || 'medium',
            likes_count: 0,
            is_liked: false,
            created_at: now,
            updated_at: now,
        };

        this.service.postRecipe(recipe, this.selectedFile).subscribe({
            next: () => {
                this.toggleCreateForm();
                this.service.getAllRecipes().subscribe(); // reload feed
                this.recipeForm.reset();
            },
            error: (err) => {
                console.error("creation failed", err)
            }
        });
    }

    // ── FILE PREVIEW ────────────────────────────────────────────
    onFileSelected(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files?.length) return;
        this.selectedFile = input.files[0];
        const objUrl = URL.createObjectURL(this.selectedFile);
        this.previewUrl.set(this.sanitizer.bypassSecurityTrustUrl(objUrl));
    }

    // ── DESCRIPTION EXPAND ──────────────────────────────────────
    toggleExpand(id: string) {
        if (this.isExpanded.has(id)) {
            this.isExpanded.delete(id);
        } else {
            this.isExpanded.add(id);
        }
    }

    // ── LIKE STUB ───────────────────────────────────────────────
    toggleLike(recipeId: string) {
        const recipe = this.query.getEntity(recipeId);
        if (!recipe) return;

        const optimisticLike = !recipe.is_liked;

        // ✅ optimistic visual update
        this.store.updateLike(recipeId, {
            is_liked: optimisticLike,
            likes_count: (recipe.likes_count ?? 0) + (optimisticLike ? 1 : -1)
        });

        this.service.toggleLike(recipeId).subscribe({
            next: (res) => {
                // ✅ replace with server truth (so not double-count)
                this.store.updateLike(recipeId, {
                    is_liked: res.liked,
                    likes_count: res.likes_count
                });
            },
            error: () => {
                // revert on error
                this.store.updateLike(recipeId, {
                    is_liked: !optimisticLike,
                    likes_count: (recipe.likes_count ?? 0)
                });
            }
        });
    }
}