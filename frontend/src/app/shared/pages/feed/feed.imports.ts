import { AsyncPipe, DatePipe } from "@angular/common";
import { Type } from "@angular/core";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { RouterLink } from "@angular/router";

export const FEEDCOMPONENTS: Type<any>[] = [
    ReactiveFormsModule,
    FormsModule,
    AsyncPipe,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatIconModule,
    DatePipe,
    RouterLink
]
