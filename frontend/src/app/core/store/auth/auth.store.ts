import { Injectable } from "@angular/core";
import { EntityState, EntityStore, StoreConfig } from "@datorama/akita";
import { Login } from "./auth.model";

export interface AuthState extends EntityState<Login> { }

@Injectable({ providedIn: 'root' })
@StoreConfig({ name: 'auth' })
export class AuthStore extends EntityStore<AuthState> {
    constructor() {
        super();
    }
}