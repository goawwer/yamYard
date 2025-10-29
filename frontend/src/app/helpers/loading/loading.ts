import { AsyncPipe } from "@angular/common";
import { Component, Input, OnInit } from "@angular/core";
import { Observable, tap } from "rxjs";
import { LoadingService } from "./loading.service";
import { RouteConfigLoadEnd, RouteConfigLoadStart, Router } from "@angular/router";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

@Component({
    selector: "loading-indicator",
    templateUrl: "./loading.html",
    styleUrls: ["./loading.css"],
    imports: [MatProgressSpinnerModule, AsyncPipe],
    standalone: true,
})
export class LoadingIndicatorComponent implements OnInit {
    loading$: Observable<boolean>;

    @Input()
    detectRouteTransitions = false;

    constructor(
        private loadingService: LoadingService,
        private router: Router
    ) {
        this.loading$ = this.loadingService.loading$;
    }

    ngOnInit(): void {
        if (this.detectRouteTransitions) {
            this.router.events.pipe(
                tap((event) => {
                    if (event instanceof RouteConfigLoadStart) {
                        this.loadingService.start();
                    } else if (event instanceof RouteConfigLoadEnd) {
                        this.loadingService.stop();
                    }
                })
            ).subscribe();
        }
    }
}