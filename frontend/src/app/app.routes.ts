import { Routes } from '@angular/router';
// import { AuthorizedPage } from './shared/pages/authorized-page/authorized-page';
// import { LoginForm } from './shared/pages/login/login-form';
import { AuthGuard } from './core/guards/auth.guard';
// import { UsersTableComponent } from './shared/pages/users-table/users-table';
// import { UserComponent } from './shared/pages/user/user';
import { HomeComponent } from './shared/pages/home/home';
import { LoginComponent } from './shared/pages/auth/login/login';
import { SignupComponent } from './shared/pages/auth/signup/signup';
import { AuthLayoutComponent } from './shared/pages/auth/authLayout/auth-layout';
import { FeedComponent } from './shared/pages/feed/feed';
import { UnauthorizedLayoutComponent } from './shared/layout/unauthorized/unauthorized-layout';
import { AuthorizedLayoutComponent } from './shared/layout/authorized/authorized-layout';
import { ProfileComponent } from './shared/pages/profile/profile';
import { AdminGuard } from './core/guards/admin.guard';
import { AdminComponent } from './shared/pages/admin/admin';
import { InfoComponent } from './shared/pages/info/info';

export const routes: Routes = [
    // Public pages
    {
        path: '',
        component: UnauthorizedLayoutComponent,
        children: [
            { path: '', redirectTo: '/home', pathMatch: 'full' },
            { path: 'home', component: HomeComponent },
            { path: 'info', component: InfoComponent }
        ]
    },

    // Auth layout (login/register)
    {
        path: 'auth',
        component: AuthLayoutComponent,
        children: [
            { path: 'login', component: LoginComponent },
            { path: 'register', component: SignupComponent },
        ]
    },

    // Main layout (protected)
    {
        path: '',
        component: AuthorizedLayoutComponent,
        canActivate: [AuthGuard],
        children: [
            // Redirect empty path '' to 'feed' for authorized users
            { path: '', redirectTo: 'feed', pathMatch: 'full' },
            { path: 'feed', component: FeedComponent },
            { path: 'me', component: ProfileComponent },
            { path: 'user/:id', component: ProfileComponent }
        ]
    },

    {
        path: 'admin',
        component: AuthorizedLayoutComponent,
        canActivate: [AuthGuard, AdminGuard],
        children: [
            { path: '', component: AdminComponent }
        ]
    },

    { path: '**', redirectTo: 'home', pathMatch: 'full' }
];

