import { Component, OnInit, signal } from '@angular/core';
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
import { debounceTime, distinctUntilChanged, map, startWith, switchMap } from 'rxjs/operators';
import { User } from '../../../core/store/user/user.model';
import { Recipe } from '../../../core/store/recipe/recipe.model';
import { AdminService } from '../../../core/store/admin/admin.service';
import { ConfirmDialogComponent } from '../../../helpers/dialog/dialog.component';
import { MatDialog } from '@angular/material/dialog';
import { UserService } from '../../../core/store/user/user.service';
import { RecipeService } from '../../../core/store/recipe/recipe.service';
import { RecipeQuery } from '../../../core/store/recipe/recipe.query';
import { UserQuery } from '../../../core/store/user/user.query';
import { MatInputModule } from '@angular/material/input';

@Component({
    selector: 'app-admin',
    standalone: true,
    imports: [
        CommonModule, ReactiveFormsModule,
        MatTableModule, MatButtonModule, MatIconModule,
        MatButtonToggleModule, MatFormFieldModule, MatSelectModule,
        MatProgressSpinnerModule, MatInputModule
    ],
    templateUrl: './admin.html',
    styleUrl: './admin.scss'
})
export class AdminComponent implements OnInit {
    mode = signal<'users' | 'recipes'>('users');
    searchCtrl = new FormControl('', { nonNullable: true });
    filteredData$!: Observable<any[]>;

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
        private dialog: MatDialog,
        private userService: UserService,
        private userQuery: UserQuery,
        private recipeService: RecipeService,
        private recipeQuery: RecipeQuery
    ) { }

    ngOnInit() {
        // Реактивный поиск с дебаунсом
        this.searchCtrl.valueChanges.pipe(
            debounceTime(300),
            distinctUntilChanged(),
            startWith('')
        ).subscribe(searchTerm => {
            this.applyFilter(searchTerm.trim());
        });

        this.loadData(); // первый загруз
    }

    switchMode(mode: 'users' | 'recipes') {
        this.mode.set(mode);
        this.searchCtrl.setValue(''); // сброс поиска
        this.loadData();
    }

    private applyFilter(searchTerm: string) {
        if (this.mode() === 'users') {
            this.userService.getAllUsers(undefined, searchTerm ? { column: 'username', value: searchTerm } : undefined)
                .subscribe();
            this.filteredData$ = this.userQuery.selectAll();
        } else {
            this.recipeService.getAllRecipes(undefined, searchTerm ? { column: 'username', value: searchTerm } : undefined)
                .subscribe();
            this.filteredData$ = this.recipeQuery.selectRecipes();
        }
    }

    private loadData() {
        this.applyFilter(this.searchCtrl.value);
    }

    clearSearch() {
        this.searchCtrl.setValue('');
    }

    deleteUser(user: User) {
        const dialogRef = this.dialog.open(ConfirmDialogComponent, {
            data: { title: 'Удалить пользователя?', message: `Удалить ${user.username}?` }
        });

        dialogRef.afterClosed().subscribe(confirmed => {
            if (confirmed) {
                this.adminService.deleteUser(user.id).subscribe();
            }
        });
    }

    deleteRecipe(recipe: Recipe) {
        const dialogRef = this.dialog.open(ConfirmDialogComponent, {
            data: { title: 'Удалить рецепт?', message: `Удалить "${recipe.title}"?` }
        });

        dialogRef.afterClosed().subscribe(confirmed => {
            if (confirmed) {
                this.adminService.deleteRecipe(recipe.id).subscribe();
            }
        });
    }

    get displayedColumns(): string[] {
        return (this.mode() === 'users' ? this.userHeaders : this.recipeHeaders).map(h => h.key);
    }
}