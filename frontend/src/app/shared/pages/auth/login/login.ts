import { Component, OnInit, signal } from '@angular/core';
import { FormControl } from '@angular/forms';
import { FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AUTHCOMPONENTS } from '../auth.imports';
import { AuthService } from '../../../../core/store/auth/auth.service';

@Component({
    selector: 'app-login-form',
    imports: [AUTHCOMPONENTS],
    templateUrl: './login.html',
    styleUrl: './login.scss'
})
export class LoginComponent implements OnInit {
    loginForm!: FormGroup;
    isSubmitted = false;

    constructor(
        private service: AuthService,
        private router: Router
    ) { }

    ngOnInit() {
        this.loginForm = new FormGroup({
            email: new FormControl('', [Validators.required, Validators.email]),
            password: new FormControl('', [Validators.required, Validators.minLength(4)])
        });
    }

    hide = signal(true);
    togglePasswordVisibility() {
        this.hide.set(!this.hide());
    }

    fields = [
        { key: 'email', label: 'Почта' },
        { key: 'password', label: 'Пароль', type: 'password' }
    ];

    signIn() {
        this.isSubmitted = true;


        this.service.signIn(this.loginForm.value).subscribe({
            next: () => {
                this.router.navigate(['/feed']);
            },
            error: (err) => {
                alert(`Login failed: ${err.message}`);
                this.isSubmitted = false;
            }
        })

    }
}
