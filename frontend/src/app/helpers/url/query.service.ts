import { Injectable } from '@angular/core';
import { Sort } from '@angular/material/sort';

@Injectable({
    providedIn: 'root'
})
export class UrlQueryService {
    constructor() { }

    filterAndSortingQuery(baseUrl: string, sort?: Sort, filter?: { column: string; value: string }): string {
        let query = baseUrl;
        let params: string[] = [];

        if (filter?.column && filter?.value) {
            params.push(`${filter.column}=${filter.value}`);
        }

        if (sort?.direction) {
            params.push(`sort=${sort.active}`, `order=${sort.direction}`);
        }

        if (params.length > 0) {
            query += `?${params.join("&")}`;
        }

        return query
    }
}