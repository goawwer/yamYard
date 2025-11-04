import { AsyncPipe, DatePipe } from "@angular/common";
import { Type } from "@angular/core";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatDividerModule } from "@angular/material/divider";
import { MatFormFieldControl, MatFormFieldModule, MatLabel } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatMenuModule } from "@angular/material/menu";
import { MatSelectModule } from "@angular/material/select";
import { MatTabsModule } from "@angular/material/tabs";
import { RouterLink } from "@angular/router";
import { MatDialogModule } from "@angular/material/dialog";
import { PostCardComponent } from "../../post-card/post-card";

export const PROFILECOMPONENTS: Type<any>[] = [
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
    RouterLink,
    MatDividerModule,
    MatLabel,
    MatTabsModule,
    MatMenuModule,
    MatDialogModule,
    PostCardComponent
]
