import { Component } from '@angular/core';
import { HeaderComponent } from './components/header-component/header.component';
import { FooterComponent } from './components/footer-component/footer.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    HeaderComponent,
    FooterComponent,
    DashboardComponent
  ],
  template: `
    <app-header-component></app-header-component>
    <app-dashboard></app-dashboard>
    <app-footer></app-footer>
  `,
})
export class App {}
