import { AsyncPipe } from "@angular/common";
import { Type } from "@angular/core";
import { MatIcon } from "@angular/material/icon";
import { RouterModule } from "@angular/router";

export const LAYOUTCOMPONENTS: Type<any>[] = [
    RouterModule,
    MatIcon,
    AsyncPipe
]