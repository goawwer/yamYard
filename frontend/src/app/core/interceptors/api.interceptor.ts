import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '../store/auth/auth.service';
import { catchError, switchMap, throwError } from 'rxjs';
import { baseURL } from '../../../environments/environment';

export const apiInterceptor: HttpInterceptorFn = (req, next) => {
    const authService = inject(AuthService);

    if (req.url.startsWith('/assets/')) return next(req);

    const apiReq = req.clone({
        headers: req.headers.set('Content-Type', 'application/json'),
        url: `${baseURL}${req.url}`,
        withCredentials: true
    });

    return next(apiReq).pipe(
        catchError((error: HttpErrorResponse) => {
            // only handle 401 errors, skip if it's refresh itself to avoid loops
            if (error.status === 401 && !req.url.includes('/refresh')) {
                return authService.refresh().pipe(
                    switchMap(success => {
                        if (success) {
                            // retry original request
                            return next(apiReq);
                        } else {
                            // refresh failed → logout
                            authService.logout().subscribe();
                            return throwError(() => error);
                        }
                    })
                );
            }

            return throwError(() => error);
        })
    );
};