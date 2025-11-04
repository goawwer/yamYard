import { Injectable } from "@angular/core";
import { RecipeStore } from "./recipe.store";
import { HttpClient, HttpContext } from "@angular/common/http";
import { Router } from "@angular/router";
import { UrlQueryService } from "../../../helpers/url/query.service";
import { Sort } from "@angular/material/sort";
import { isUpdatingRecipe, Recipe } from "./recipe.model";
import { catchError, finalize, Observable, tap, throwError } from "rxjs";
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
}