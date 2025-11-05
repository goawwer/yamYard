import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { Observable, of, combineLatest } from 'rxjs';
import { map, startWith, switchMap } from 'rxjs/operators';
import { User } from '../../../core/store/user/user.model';
import { Recipe } from '../../../core/store/recipe/recipe.model';
import { AdminService } from '../../../core/store/admin/admin.service';
import { ConfirmDialogComponent } from '../../../helpers/dialog/dialog.component';
import { MatDialog } from '@angular/material/dialog';

@Component({
    selector: 'app-admin',
    standalone: true,
    imports: [
        CommonModule, ReactiveFormsModule,
        MatTableModule, MatButtonModule, MatIconModule,
        MatButtonToggleModule, MatFormFieldModule, MatSelectModule,
        MatProgressSpinnerModule
    ],
    templateUrl: './admin.html',
    styleUrl: './admin.scss'
})
export class AdminComponent implements OnInit {
    mode = 'users'; // 'users' | 'recipes'

    filterForm = new FormGroup({
        username: new FormControl('')
    });

    users$: Observable<User[]>;
    recipes$: Observable<Recipe[]>;
    filteredData$: Observable<any[]>;

    userHeaders = [
        { key: 'username', value: 'Имя пользователя' },
        { key: 'email', value: 'Email' },
        { key: 'is_admin', value: 'Админ' },
        { key: 'created_at', value: 'Создан' },
        { key: 'actions', value: 'Действия' }
    ];

    recipeHeaders = [
        { key: 'title', value: 'Название' },
        { key: 'author_username', value: 'Автор' },
        { key: 'likes_count', value: 'Лайки' },
        { key: 'created_at', value: 'Создан' },
        { key: 'actions', value: 'Действия' }
    ];

    constructor(
        private adminService: AdminService,
        private dialog: MatDialog
    ) {
        this.users$ = this.adminService.getUsers();
        this.recipes$ = this.adminService.getRecipes();

        this.filteredData$ = combineLatest([
            this.users$,
            this.recipes$,
            this.filterForm.get('username')!.valueChanges.pipe(startWith(''))
        ]).pipe(
            map(([users, recipes, filter]) => {
                if (this.mode === 'users') {
                    return filter ? users.filter(u => u.username.includes(filter)) : users;
                } else {
                    return filter ? recipes.filter(r => r.author_username?.includes(filter)) : recipes;
                }
            })
        );
    }

    ngOnInit() { }

    switchMode(mode: 'users' | 'recipes') {
        this.mode = mode;
        this.filterForm.get('username')?.setValue('');
    }

    deleteUser(user: User) {
        const dialogRef = this.dialog.open(ConfirmDialogComponent, {
            data: { title: 'Удалить пользователя?', message: `Удалить ${user.username}?` }
        });

        dialogRef.afterClosed().subscribe(confirmed => {
            if (confirmed) {
                this.adminService.deleteUser(user.id).subscribe(() => {
                    this.users$ = this.adminService.getUsers();
                });
            }
        });
    }

    deleteRecipe(recipe: Recipe) {
        const dialogRef = this.dialog.open(ConfirmDialogComponent, {
            data: { title: 'Удалить рецепт?', message: `Удалить "${recipe.title}"?` }
        });

        dialogRef.afterClosed().subscribe(confirmed => {
            if (confirmed) {
                this.adminService.deleteRecipe(recipe.id).subscribe(() => {
                    this.recipes$ = this.adminService.getRecipes();
                });
            }
        });
    }

    get displayedColumns(): string[] {
        return (this.mode === 'users' ? this.userHeaders : this.recipeHeaders).map(h => h.key);
    }
}