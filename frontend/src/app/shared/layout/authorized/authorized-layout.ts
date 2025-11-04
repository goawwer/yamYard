import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../../core/store/auth/auth.service';
import { LAYOUTCOMPONENTS } from '../layout.imports';
import { UserService } from '../../../core/store/user/user.service';
import { Observable } from 'rxjs';
import { User } from '../../../core/store/user/user.model';
import { UserQuery } from '../../../core/store/user/user.query';

@Component({
    selector: 'app-layout',
    templateUrl: './authorized-layout.html',
    styleUrls: ['./authorized-layout.scss'],
    imports: [LAYOUTCOMPONENTS]
})
export class AuthorizedLayoutComponent implements OnInit {
    currentYear = new Date().getFullYear();
    user$!: Observable<User | undefined>;

    constructor(
        private userService: UserService,
        private service: AuthService,
        private router: Router
    ) {
    }

    ngOnInit(): void {
        this.user$ = this.userService.getAuthenticatedUser();
    }

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
        return this.router.url.includes("/me")
    }
}
