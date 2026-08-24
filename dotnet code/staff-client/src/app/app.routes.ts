import { Routes } from '@angular/router';

import { SearchResultComponent } from './features/staff-directory/searchresult.component';

export const routes: Routes = [
  { path: '', component: SearchResultComponent },
  { path: '**', redirectTo: '' }
];
