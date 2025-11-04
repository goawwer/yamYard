import { Component, OnInit, signal } from '@angular/core';
import { UserService } from '../../../core/store/user/user.service';
import { combineLatest, EMPTY, filter, finalize, map, Observable, of, switchMap, take } from 'rxjs';
import { isUpdatingUser, User } from '../../../core/store/user/user.model';
import { UserQuery } from '../../../core/store/user/user.query';
import { RecipeService } from '../../../core/store/recipe/recipe.service';
import { ActivatedRoute } from '@angular/router';
import { Sort } from '@angular/material/sort';
import { isUpdatingRecipe, Recipe } from '../../../core/store/recipe/recipe.model';
import { FormControl, FormGroup } from '@angular/forms';
import { PROFILECOMPONENTS } from './profile.imports';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { MatDialog } from '@angular/material/dialog';
import { ConfirmDialogComponent } from '../../../helpers/dialog/dialog.component';

@Component({
    selector: 'app-profile',
    imports: [PROFILECOMPONENTS],
    templateUrl: './profile.html',
    styleUrl: './profile.scss',
})
export class ProfileComponent implements OnInit {
    isUpdating = signal(false);
    user$!: Observable<User | undefined>;
    recipes$!: Observable<Recipe[]>;
    editForm = new FormGroup({
        username: new FormControl(''),
        profileStatus: new FormControl(''),
        bio: new FormControl(''),
    });
    selectedFile?: File;
    previewUrl = signal<SafeUrl | null>(null);
    currentAvatarUrl = signal<string | null>(null);
    isExpanded = new Set<string>();
    isOwnProfile = signal(false);

    // Add to signals
    editingRecipeId = signal<string | null>(null);
    isEditingRecipe = signal(false);
    recipePreviewUrl = signal<SafeUrl | null>(null);
    currentRecipeImage = signal<string | null>(null);

    // New form
    recipeEditForm = new FormGroup({
        title: new FormControl('', { nonNullable: true }),
        description: new FormControl('', { nonNullable: true }),
        ingredients: new FormControl(''),
        cooking_time: new FormControl(0),
        difficulty: new FormControl('', { nonNullable: true }),
        image_url: new FormControl('', { nonNullable: true }),
    });

    selectedRecipeFile?: File;

    constructor(
        private userService: UserService,
        private userQuery: UserQuery,
        private recipeService: RecipeService,
        private route: ActivatedRoute,
        private sanitizer: DomSanitizer,
        private dialog: MatDialog
    ) { }

    ngOnInit(): void {
        const isMeRoute = this.route.snapshot.routeConfig?.path === 'me';

        this.user$ = isMeRoute
            ? this.userService.getAuthenticatedUser()
            : this.route.paramMap.pipe(
                switchMap(params => {
                    const id = params.get('id');
                    if (!id) return EMPTY;
                    return this.userService.getUserById(id);
                })
            );

        this.recipes$ = this.user$.pipe(
            switchMap(user => {
                if (!user?.username) return of([]);
                return this.recipeService.getAllRecipes(
                    { active: 'created_at', direction: 'desc' } as Sort,
                    { column: 'username', value: user.username }
                );
            })
        );

        combineLatest([this.user$, this.userQuery.selectAuthenticated()])
            .subscribe(([profileUser, authUser]) => {
                this.isOwnProfile.set(!!authUser && profileUser?.id === authUser.id);
                this.currentAvatarUrl.set(profileUser?.image_url ?? null);
            });
    }

    onFileSelected(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files?.length) return;

        this.selectedFile = input.files[0];

        // create a safe preview URL
        const objectUrl = URL.createObjectURL(this.selectedFile);
        this.previewUrl.set(this.sanitizer.bypassSecurityTrustUrl(objectUrl));
    }

    editProfile() {
        const userValues: Partial<isUpdatingUser> = this.editForm.value;

        this.user$.pipe(take(1)).subscribe(activeUser => {
            if (!activeUser) return;

            this.userService.updateUser(activeUser.id, userValues, this.selectedFile)
                .subscribe({
                    next: () => {
                        this.isUpdating.set(false);
                        this.editForm.reset();
                        this.selectedFile = undefined;
                    },
                    error: (err) => {
                        console.error('Update failed', err);
                    }
                });
        });
    }

    toggleUpdateForm() {
        this.isUpdating.update(v => !v);

        if (this.isUpdating()) {
            this.user$.pipe(take(1)).subscribe(user => {
                if (user) {
                    this.editForm.patchValue({
                        username: user.username,
                        bio: user.bio,
                        profileStatus: user.profileStatus,
                    });
                }
            });
        }
    }

    toggleExpand(id: string) {
        if (this.isExpanded.has(id)) {
            this.isExpanded.delete(id);
        } else {
            this.isExpanded.add(id);
        }
    }

    onOverlayClick(event: MouseEvent) {
        const target = event.target as HTMLElement;

        if (!target.classList.contains('overlay')) return;

        if (this.isUpdating()) {
            this.toggleUpdateForm();
        } else if (this.isEditingRecipe()) {
            this.closeRecipeEdit();
        }
    }


    openRecipeEdit(recipe: Recipe) {
        this.editingRecipeId.set(recipe.id);
        this.currentRecipeImage.set(recipe.image_url ?? null);

        this.recipeEditForm.patchValue({
            title: recipe.title,
            description: recipe.description,
            ingredients: recipe.ingredients,
            cooking_time: recipe.cooking_time,
            difficulty: recipe.difficulty,
            image_url: recipe.image_url,
        });

        this.recipePreviewUrl.set(null);
        this.selectedRecipeFile = undefined;
        this.isEditingRecipe.set(true);
    }

    closeRecipeEdit() {
        this.isEditingRecipe.set(false);
        this.editingRecipeId.set(null);
        this.recipePreviewUrl.set(null);
        this.selectedRecipeFile = undefined;
        this.recipeEditForm.reset();
    }

    onRecipeFileSelected(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files?.length) return;

        this.selectedRecipeFile = input.files[0];
        const objectUrl = URL.createObjectURL(this.selectedRecipeFile);
        this.recipePreviewUrl.set(this.sanitizer.bypassSecurityTrustUrl(objectUrl));
    }

    updateRecipe() {
        if (!this.editingRecipeId() || this.recipeEditForm.invalid) return;

        const formValue: Partial<isUpdatingRecipe> = this.recipeEditForm.getRawValue();

        this.recipeService.updateRecipe(
            this.editingRecipeId()!,
            formValue,
            this.selectedRecipeFile
        ).pipe(
            switchMap(() => this.recipeService.getAllRecipes()),
            take(1)
        ).subscribe({
            next: (updatedRecipes) => {
                this.closeRecipeEdit();
                this.recipes$ = of(updatedRecipes);
            },
            error: (err) => console.error('Update recipe failed', err)
        });
    }

    deleteRecipe(recipe: Recipe) {
        const dialogRef = this.dialog.open(ConfirmDialogComponent, {
            width: '320px',
            data: {
                title: 'Удалить рецепт?',
                message: `Вы уверены, что хотите удалить "${recipe.title}"?`,
            },
        });

        dialogRef.afterClosed().subscribe((confirmed) => {
            if (confirmed) {
                this.recipeService.deleteRecipe(recipe.id).pipe(
                    switchMap(() => this.recipeService.getAllRecipes()),
                    finalize(() => this.dialog.closeAll())
                ).subscribe((recipes) => {
                    this.recipes$ = of(recipes);
                });
            }
        });
    }

    toggleLike(id: string) {
        console.log('Toggle like', id);
    }
}