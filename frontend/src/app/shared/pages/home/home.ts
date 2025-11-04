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
}
