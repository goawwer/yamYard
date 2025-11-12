import { Component, Input, Output, EventEmitter, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule, AsyncPipe } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { RouterLink } from '@angular/router';
import { Recipe } from '../../core/store/recipe/recipe.model';
import { Observable } from 'rxjs';
import { RecipeQuery } from '../../core/store/recipe/recipe.query';

@Component({
    selector: 'app-post-card',
    standalone: true,
    imports: [
        CommonModule,
        AsyncPipe,
        MatCardModule,
        MatIconModule,
        MatButtonModule,
        MatMenuModule,
        RouterLink
    ],
    templateUrl: './post-card.html',
    styleUrl: './post-card.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class PostCardComponent {
    @Input({ required: true }) recipe!: Recipe;
    @Input() canManage = false;
    @Input() avatarURL?: string | null = null;

    @Output() edit = new EventEmitter<void>();
    @Output() delete = new EventEmitter<void>();
    @Output() like = new EventEmitter<string>();

    recipe$!: Observable<Recipe | undefined>;

    expanded = signal(false);

    constructor(private query: RecipeQuery) { }

    ngOnInit() {
        this.recipe$ = this.query.selectEntity(this.recipe.id);
    }

    toggleLike() {
        this.like.emit(this.recipe.id);
    }

    toggleExpand() {
        this.expanded.update(v => !v);
    }

    onEdit() {
        this.edit.emit();
    }

    onDelete() {
        this.delete.emit();
    }

    translateDifficulty(diff?: string): string {
        if (!diff) return '';
        const map: Record<string, string> = {
            easy: 'Легко',
            medium: 'Средне',
            hard: 'Сложно'
        };
        return map[diff.toLowerCase()] ?? diff;
    }

    formatYekaterinburg(dateStr: string): string {
        const date = new Date(dateStr);
        return date.toLocaleString('ru-RU', { timeZone: 'Asia/Yekaterinburg' });
    }
}