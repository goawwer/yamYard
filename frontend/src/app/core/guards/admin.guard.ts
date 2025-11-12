import { Injectable } from "@angular/core";
import { CanActivate, Router } from "@angular/router";
import { UserQuery } from "../store/user/user.query";
import { filter, map, Observable, take } from "rxjs";
import { User } from "../store/user/user.model";

@Injectable({ providedIn: 'root' })
export class AdminGuard implements CanActivate {
    constructor(private userQuery: UserQuery, private router: Router) { }

    canActivate(): Observable<boolean> {
        return this.userQuery.selectAuthenticated().pipe(
            // wait until user is non-null
            filter((user): user is User => user !== undefined && user !== null),
            take(1),
            map(user => {
                if (!user.is_admin) {
                    this.router.navigate(['/']);
                    return false;
                }
                return true;
            })
        );
    }
}