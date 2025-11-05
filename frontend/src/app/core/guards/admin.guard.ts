import { Injectable } from "@angular/core";
import { CanActivate, Router } from "@angular/router";
import { UserQuery } from "../store/user/user.query";
import { map, Observable } from "rxjs";

@Injectable({ providedIn: 'root' })
export class AdminGuard implements CanActivate {
    constructor(private userQuery: UserQuery, private router: Router) { }

    canActivate(): Observable<boolean> {
        return this.userQuery.selectAuthenticated().pipe(
            map(user => {
                if (!user?.is_admin) {
                    this.router.navigate(['/']);
                    return false;
                }
                return true;
            })
        );
    }
}