import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { BehaviorSubject, catchError, Observable, of, tap } from "rxjs";
import { User } from "../user/user.model";
import { Recipe } from "../recipe/recipe.model";

@Injectable({ providedIn: 'root' })
export class AdminService {
    private adminURL = '/api/admin';
    private usersSubject = new BehaviorSubject<User[]>([]);
    private recipesSubject = new BehaviorSubject<Recipe[]>([]);

    users$ = this.usersSubject.asObservable();
    recipes$ = this.recipesSubject.asObservable();

    constructor(private http: HttpClient) { }

    loadUsers() {
        this.http.get<User[]>(`${this.adminURL}/users`).pipe(
            catchError(() => of([]))
        ).subscribe(users => this.usersSubject.next(users));
    }

    loadRecipes() {
        this.http.get<Recipe[]>(`${this.adminURL}/recipes`).pipe(
            catchError(() => of([]))
        ).subscribe(recipes => this.recipesSubject.next(recipes));
    }

    deleteUser(id: string): Observable<void> {
        return this.http.delete<void>(`${this.adminURL}/users/${id}`).pipe(
            tap(() => this.loadUsers())
        );
    }

    deleteRecipe(id: string): Observable<void> {
        return this.http.delete<void>(`${this.adminURL}/recipes/${id}`).pipe(
            tap(() => this.loadRecipes())
        );
    }
}