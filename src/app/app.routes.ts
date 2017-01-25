import { Routes } from '@angular/router';

import { HomeComponent } from './home/home.component';
import { ContactComponent} from './contact/contact.component';
import { PrivacyComponent} from './privacy/privacy.component';
import { FailureComponent } from './failure/failure.component';

export const rootRouterConfig: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', component: HomeComponent },
  { path: 'contact', component: ContactComponent },
  { path: 'privacy', component: PrivacyComponent },
  { path: 'failure', component: FailureComponent }
];

