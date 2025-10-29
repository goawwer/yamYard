import { Injectable } from "@angular/core";
import { UserStore } from "./user.store";
import { HttpClient } from "@angular/common/http";
import { User } from "./user.model";
import { catchError, finalize, Observable, tap, throwError } from "rxjs";
import { Sort } from "@angular/material/sort";
import { UrlQueryService } from "../../../helpers/url/query.service";
import { Router } from "@angular/router";

@Injectable({
    providedIn: 'root'
})
export class UserService {
    private getAllUsersURL = '/api/users';

    constructor(
        readonly store: UserStore,
        private http: HttpClient,
        private router: Router,
        private urlHelper: UrlQueryService
    ) { }

    getAllUsers(sort?: Sort, filter?: { column: string; value: string }): void {
        this.store.setLoading(true);
        let query = this.urlHelper.filterAndSortingQuery(this.getAllUsersURL, sort, filter);

        this.http.get<User[]>(query).pipe(
            tap(users => {
                this.store.set(users);
            }),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error || 'failed to get users'));
            }),
            finalize(() => this.store.setLoading(false))
        )
    }

    getAuthenticatedUser(): Observable<User> {
        this.store.setLoading(true);

        return this.http.get<User>(`${this.getAllUsersURL}/me`).pipe(
            catchError((err) => {
                return throwError(() => new Error(err.error?.error || 'failed to get users'));
            }),
            finalize(() => this.store.setLoading(false))
        );
    }

    getUserById(id: number): Observable<User> {
        this.store.setLoading(true);
        return this.http.get<User>(`${this.getAllUsersURL}/${id}`).pipe(
            catchError((err) => {
                this.router.navigate(['users/create'])
                return throwError(() => new Error(err.error?.error || 'user not found'))
            }),
            finalize(() => this.store.setLoading(false))
        )
    }

    postUser(user: User): Observable<User> {
        this.store.setLoading(true);
        return this.http.post<User>(`${this.getAllUsersURL}/create`, user).pipe(
            tap(created => this.store.add(created)),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error) || 'failed to create user');
            }),
            finalize(() => this.store.setLoading(false))
        );
    }

    updateUser(id: number, user: User): Observable<User> {
        this.store.setLoading(true);
        return this.http.put<User>(`${this.getAllUsersURL}/${id}/update`, user).pipe(
            tap(updated => this.store.update(updated)),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error) || 'failed to update user');
            }),
            finalize(() => this.store.setLoading(false))
        )
    }
}