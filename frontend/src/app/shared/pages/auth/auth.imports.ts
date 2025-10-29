import { Type } from "@angular/core";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatError, MatFormFieldModule } from "@angular/material/form-field";
import { MatIcon, MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { RouterLink } from "@angular/router";

export const AUTHCOMPONENTS: Type<any>[] = [
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatError,
    MatInputModule,
    MatIconModule,
    MatButtonModule,
    RouterLink,
    MatIcon
]
