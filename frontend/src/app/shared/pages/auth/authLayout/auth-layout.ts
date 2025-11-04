import { Component } from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import { RouterModule } from '@angular/router';

@Component({
    selector: 'app-layout',
    templateUrl: './auth-layout.html',
    imports: [RouterModule, MatIcon]
})
export class AuthLayoutComponent {
    currentYear = new Date().getFullYear();
}
