import { HttpContextToken, HttpInterceptorFn } from "@angular/common/http";
import { inject } from "@angular/core";
import { LoadingService } from "../../helpers/loading/loading.service";
import { finalize } from "rxjs";

export const skipLoading = new HttpContextToken<boolean>(() => false);

export const loadingInterceptor: HttpInterceptorFn = (req, next) => {
    const loadingService = inject(LoadingService);

    if (req.context.get(skipLoading)) {
        return next(req);
    }

    loadingService.start();

    return next(req).pipe(
        finalize(() => { loadingService.stop() })
    )
}
