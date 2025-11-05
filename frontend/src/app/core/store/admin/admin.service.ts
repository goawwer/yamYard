import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { catchError, Observable, of } from "rxjs";
import { User } from "../user/user.model";
import { Recipe } from "../recipe/recipe.model";

@Injectable({ providedIn: 'root' })
export class AdminService {
    private adminURL = '/adminURL/admin';

    constructor(private http: HttpClient) { }

    getUsers(): Observable<User[]> {
        return this.http.get<User[]>(`${this.adminURL}/users`).pipe(
            catchError(() => of([]))
        );
    }

    deleteUser(id: string): Observable<void> {
        return this.http.delete<void>(`${this.adminURL}/users/${id}`);
    }

    getRecipes(): Observable<Recipe[]> {
        return this.http.get<Recipe[]>(`${this.adminURL}/recipes`).pipe(
            catchError(() => of([]))
        );
    }

    deleteRecipe(id: string): Observable<void> {
        return this.http.delete<void>(`${this.adminURL}/recipes/${id}`);
    }
}