import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { LAYOUTCOMPONENTS } from '../layout.imports';

@Component({
    selector: 'app-layout',
    templateUrl: './unauthorized-layout.html',
    styleUrls: ['./unauthorized-layout.scss'],
    imports: [LAYOUTCOMPONENTS]
})
export class UnauthorizedLayoutComponent {
    currentYear = new Date().getFullYear();
}
