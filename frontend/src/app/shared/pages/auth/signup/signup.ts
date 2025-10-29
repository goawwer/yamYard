import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { FormGroup, Validators } from '@angular/forms';
import { AUTHCOMPONENTS } from '../auth.imports';

@Component({
    selector: 'app-signup',
    imports: [AUTHCOMPONENTS],
    templateUrl: './signup.html',
    styleUrl: './signup.scss'
})
export class SignupComponent implements OnInit {
    signupForm!: FormGroup;
    isSubmitted = false;
    hidePassword = true;

    constructor() {
    }

    ngOnInit() {
        this.signupForm = new FormGroup({
            username: new FormControl('', [Validators.required, Validators.minLength(5)]),
            email: new FormControl('', [Validators.required, Validators.email]),
            password: new FormControl('', [Validators.required, Validators.minLength(8)]),
            description: new FormControl('',)
        });
    }

    togglePasswordVisibility() {
        this.hidePassword = !this.hidePassword;
    }

    signUp() {
        if (this.signupForm.invalid) return;
        this.isSubmitted = true;

        const user = this.signupForm.value;
        console.log('Registering user:', user);

        // TODO: call API endpoint
    }
}
