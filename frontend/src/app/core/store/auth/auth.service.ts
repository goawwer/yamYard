import { BehaviorSubject, catchError, finalize, map, Observable, of, switchMap, tap } from "rxjs";
import { AuthStore } from "./auth.store";
import { HttpClient } from "@angular/common/http";
import { UserService } from "../user/user.service";
import { Login } from "./auth.model";
import { Injectable } from "@angular/core";

@Injectable({ providedIn: 'root' })
export class AuthService {
    private loginUrl = '/auth/login';
    private checkUrl = '/api/check';
    private refreshUrl = '/auth/refresh';
    private logoutUrl = '/api/logout';

    private isAuthenticatedSubject = new BehaviorSubject<boolean>(false);

    constructor(
        private store: AuthStore,
        private http: HttpClient,
        private userService: UserService
    ) { }

    signIn(user: Login): Observable<void> {
        this.store.setLoading(true);
        return this.http.post(this.loginUrl, user).pipe(
            switchMap(() => this.loadAuthenticatedUser()),
            finalize(() => this.store.setLoading(false))
        );
    }

    check(): Observable<boolean> {
        return this.http.get(this.checkUrl).pipe(
            switchMap(() => this.loadAuthenticatedUser()),
            map(() => {
                this.isAuthenticatedSubject.next(true);
                return true;
            }),
            catchError(() => {
                this.isAuthenticatedSubject.next(false);
                return of(false);
            })
        );
    }

    refresh(): Observable<boolean> {
        return this.http.post(this.refreshUrl, {}).pipe(
            switchMap(() => this.loadAuthenticatedUser()),
            map(() => true),
            catchError(() => of(false))
        );
    }

    logout(): Observable<void> {
        this.store.setLoading(true);
        return this.http.post<void>(this.logoutUrl, {}).pipe(
            tap(() => {
                this.isAuthenticatedSubject.next(false);
                this.userService.store.remove();
            }),
            finalize(() => this.store.setLoading(false))
        );
    }

    private loadAuthenticatedUser(): Observable<void> {
        return this.userService.getAuthenticatedUser().pipe(
            tap(user => {
                this.userService.store.set([user]);
                this.userService.store.setActive(user.id);
                this.isAuthenticatedSubject.next(true);
            }),
            map(() => void 0)
        );
    }
}
