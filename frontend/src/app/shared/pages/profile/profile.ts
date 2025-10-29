import { Component, OnInit } from '@angular/core';
import { AsyncPipe, CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIcon, MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { UserService } from '../../../core/store/user/user.service';
import { Observable } from 'rxjs';
import { User } from '../../../core/store/user/user.model';
import { UserQuery } from '../../../core/store/user/user.query';

@Component({
    selector: 'app-profile',
    imports: [CommonModule, MatCardModule, MatButtonModule, MatIconModule, MatDividerModule, AsyncPipe, MatIcon],
    templateUrl: './profile.html',
    styleUrl: './profile.scss',
})
export class ProfileComponent implements OnInit {
    user$!: Observable<User | undefined>;

    constructor(private service: UserService, private query: UserQuery) { }

    ngOnInit(): void {
        this.service.getAuthenticatedUser().subscribe();
        this.user$ = this.query.selectAuthenticated();
    }
}
