import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
    selector: 'app-home',
    templateUrl: './home.html',
    styleUrls: ['./home.scss'],
    imports: [RouterLink],
})
export class HomeComponent {
    currentYear = new Date().getFullYear();

    demoRecipes = [
        {
            title: 'Homemade Pizza',
            description: 'Crispy crust, melted cheese — pure joy in every bite.',
            image: 'assets/demo/pizza.jpg',
        },
        {
            title: 'Avocado Toast',
            description: 'Simple, quick, and always Instagram-worthy.',
            image: 'assets/demo/avocado.jpg',
        },
        {
            title: 'Berry Smoothie',
            description: 'Fresh, sweet, and full of morning energy.',
            image: 'assets/demo/smoothie.jpg',
        },
    ];
}
