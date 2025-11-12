import { Component, signal } from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { RouterLink } from '@angular/router';
import { animate, style, transition, trigger } from '@angular/animations';

interface FaqItem {
    question: string;
    answer: string;
    expanded: boolean;
}

@Component({
    selector: 'app-info',
    templateUrl: './info.html',
    styleUrls: ['./info.scss'],
    imports: [MatIcon, MatCardModule, RouterLink],
    animations: [
        trigger('fadeSlide', [
            transition(':enter', [
                style({ opacity: 0, height: 0, marginBottom: 0 }),
                animate('300ms ease-out', style({ opacity: 1, height: '*', marginBottom: '1rem' }))
            ]),
            transition(':leave', [
                animate('250ms ease-in', style({ opacity: 0, height: 0, marginBottom: 0 }))
            ])
        ])
    ]
})
export class InfoComponent {
    faqItems = signal<FaqItem[]>([
        {
            question: 'Что такое YamYard?',
            answer: `Это сообщество людей, которые любят готовить и делиться рецептами.<br>
               Здесь вы можете сохранять свои любимые блюда, искать новые идеи и комментировать рецепты других пользователей.`,
            expanded: false
        },
        {
            question: 'Нужно ли регистрироваться?',
            answer: 'Да, чтобы сохранять и публиковать рецепты, а также ставить оценки, потребуется создать учетную запись. Это займет не больше минуты.',
            expanded: false
        },
        {
            question: 'Можно ли добавлять фото?',
            answer: 'Конечно! YamYard поддерживает загрузку фотографий, чтобы рецепты выглядели живыми и понятными.',
            expanded: false
        },
        {
            question: 'YamYard бесплатный?',
            answer: 'Да, все основные функции проекта останутся бесплатными. Возможно, позже появятся дополнительные опции для активных пользователей.',
            expanded: false
        },
        {
            question: 'Как поддержать проект?',
            answer: 'Самый простой способ — делиться YamYard с друзьями и оставлять обратную связь. Мы всегда рады идеям и предложениям!',
            expanded: false
        }
    ]);

    toggleFaq(item: FaqItem) {
        this.faqItems.update(items =>
            items.map(i => ({
                ...i,
                expanded: i === item ? !i.expanded : false
            }))
        );
    }
}