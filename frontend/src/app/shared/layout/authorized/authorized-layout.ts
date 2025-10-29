import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../../core/store/auth/auth.service';
import { LAYOUTCOMPONENTS } from '../layout.imports';

@Component({
    selector: 'app-layout',
    templateUrl: './authorized-layout.html',
    styleUrls: ['./authorized-layout.scss'],
    imports: [LAYOUTCOMPONENTS]
})
export class AuthorizedLayoutComponent {
    currentYear = new Date().getFullYear();

    constructor(
        private service: AuthService,
        private router: Router
    ) { }

    logout() {
        this.service.logout().subscribe({
            next: () => {
                this.router.navigate(['/login']);
            },
            error: (err) => {
                alert(`Logout failed: ${err.message}`);
            }
        })
    }

    isProfilePage(): boolean {
        return this.router.url.includes("/profile")
    }
}
