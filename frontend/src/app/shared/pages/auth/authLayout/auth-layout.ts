import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';

@Component({
    selector: 'app-layout',
    templateUrl: './auth-layout.html',
    styleUrls: ['./auth-layout.scss'],
    imports: [RouterModule]
})
export class AuthLayoutComponent {
    currentYear = new Date().getFullYear();
}
