import { Component, Input, Output, EventEmitter, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { RouterLink } from '@angular/router';
import { Recipe } from '../../core/store/recipe/recipe.model';

@Component({
    selector: 'app-post-card',
    standalone: true,
    imports: [
        CommonModule,
        MatCardModule,
        MatIconModule,
        MatButtonModule,
        MatMenuModule,
        RouterLink
    ],
    templateUrl: './post-card.html',
    styleUrl: './post-card.scss'
})
export class PostCardComponent {
    @Input({ required: true }) recipe!: Recipe;
    @Input() canManage = false; // true for profile page
    @Input() currentAvatarUrl?: string | null = null;

    @Output() edit = new EventEmitter<void>();
    @Output() delete = new EventEmitter<void>();
    @Output() like = new EventEmitter<void>();

    // expanded description logic
    expanded = signal(false);

    toggleExpand() {
        this.expanded.set(!this.expanded());
    }

    toggleLike() {
        this.like.emit();
    }

    onEdit() {
        this.edit.emit();
    }

    onDelete() {
        this.delete.emit();
    }

    translateDifficulty(diff?: string) {
        if (!diff) return '';
        const map: Record<string, string> = {
            easy: 'Легко',
            medium: 'Средне',
            hard: 'Сложно'
        };
        return map[diff.toLowerCase()] ?? diff;
    }

    formatYekaterinburg(dateStr: string) {
        const date = new Date(dateStr);
        return date.toLocaleString('ru-RU', { timeZone: 'Asia/Yekaterinburg' });
    }
}
