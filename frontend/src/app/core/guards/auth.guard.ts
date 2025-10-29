import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { AuthService } from '../store/auth/auth.service';

@Injectable({
    providedIn: 'root'
})
export class AuthGuard implements CanActivate {
    constructor(
        private auth: AuthService,
        private router: Router
    ) { }

    canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
        return this.auth.check().pipe(
            tap(isAuth => {
                if (!isAuth && state.url !== '/home') {
                    alert("Unauthorized");
                    this.router.navigate(['/home']);
                } else if (isAuth && state.url === '/') {
                    this.router.navigate(['/feed'])
                }
                return true;
            })
        );
    }
}   
