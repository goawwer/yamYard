import { Injectable } from "@angular/core";
import { RecipeStore } from "./recipe.store";
import { HttpClient, HttpContext } from "@angular/common/http";
import { Router } from "@angular/router";
import { UrlQueryService } from "../../../helpers/url/query.service";
import { Sort } from "@angular/material/sort";
import { isUpdatingRecipe, Recipe } from "./recipe.model";
import { catchError, finalize, map, Observable, of, tap, throwError } from "rxjs";
import { skipJsonContentType } from "../../interceptors/api.interceptor";

@Injectable({
    providedIn: 'root'
})
export class RecipeService {
    private getAllRecipeURL = '/api/recipes';

    constructor(
        readonly store: RecipeStore,
        private http: HttpClient,
        private router: Router,
        private urlHelper: UrlQueryService
    ) { }

    getAllRecipes(sort?: Sort, filter?: { column: string; value: string }): Observable<Recipe[]> {
        this.store.setLoading(true);
        let query = this.urlHelper.filterAndSortingQuery(this.getAllRecipeURL, sort, filter);

        return this.http.get<Recipe[]>(query).pipe(
            tap(users => {
                this.store.set(users);
            }),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error || 'failed to get recipes'));
            }),
            finalize(() => this.store.setLoading(false))
        )
    }

    getRecipeById(id: string): Observable<Recipe> {
        this.store.setLoading(true);
        return this.http.get<Recipe>(`${this.getAllRecipeURL}/${id}`).pipe(
            catchError((err) => {
                this.router.navigate(['users/create'])
                return throwError(() => new Error(err.error?.error || 'recipe not found'))
            }),
            finalize(() => this.store.setLoading(false))
        )
    }

    postRecipe(recipe: Recipe, file?: File): Observable<Recipe> {
        this.store.setLoading(true);

        const formData = new FormData();
        if (file) formData.append('recipe', file);

        Object.entries(recipe).forEach(([key, value]) => {
            if (value != null) formData.append(key, value.toString());
        });

        return this.http.post<Recipe>(`${this.getAllRecipeURL}/create`, formData, {
            context: new HttpContext().set(skipJsonContentType, true)
        }).pipe(
            tap(created => this.store.add(created)),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error) || 'failed to create recipe');
            }),
            finalize(() => this.store.setLoading(false))
        );
    }

    updateRecipe(id: string, recipe: Partial<isUpdatingRecipe>, file?: File): Observable<Recipe> {
        this.store.setLoading(true);

        this.store.setLoading(true);

        const formData = new FormData();
        if (file) formData.append('recipe', file);

        Object.entries(recipe).forEach(([key, value]) => {
            if (value != null) formData.append(key, value.toString());
        });

        return this.http.put<Recipe>(`${this.getAllRecipeURL}/${id}/update`, formData, {
            context: new HttpContext().set(skipJsonContentType, true)
        }).pipe(
            tap(updated => this.store.update(updated)),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error) || 'failed to update recipe');
            }),
            finalize(() => this.store.setLoading(false))
        )
    }

    deleteRecipe(id: string): Observable<void> {
        this.store.setLoading(true);

        return this.http.delete<void>(`${this.getAllRecipeURL}/${id}/delete`).pipe(
            tap(() => this.store.remove(id)),
            catchError((err) => {
                return throwError(() => new Error(err.error?.error) || 'failed to delete recipe');
            }),
            finalize(() => this.store.setLoading(false))
        )
    }

    toggleLike(recipeId: string): Observable<{ liked: boolean; likes_count: number }> {
        return this.http.post<{ liked: boolean; likes_count: number }>(
            `${this.getAllRecipeURL}/${recipeId}/like`, {}
        ).pipe(
            tap(res => {
                this.store.updateLike(recipeId, {
                    likes_count: res.likes_count,
                    is_liked: res.liked
                });

                if (res.liked) {
                    this.store.addToLiked(recipeId);
                } else {
                    this.store.removeFromLiked(recipeId);
                }
            }),
            catchError(err => {
                console.error('Toggle like failed', err);
                return throwError(() => err);
            })
        );
    }

    getLikedRecipes(sort?: Sort, filter?: { column: string; value: string }): Observable<Recipe[]> {
        this.store.setLoading(true);
        const url = this.urlHelper.filterAndSortingQuery(`${this.getAllRecipeURL}/liked`, sort, filter);

        return this.http.get<Recipe[] | null>(url).pipe(
            map(recipes => recipes ?? []),
            tap(recipes => {
                this.store.setLikedRecipes(recipes);
            }),
            catchError(err => {
                console.error('Failed to load liked recipes', err);
                return of([]);
            }),
            finalize(() => this.store.setLoading(false))
        );
    }
}