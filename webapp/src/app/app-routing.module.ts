import { NgModule } from '@angular/core';
import { ActivatedRouteSnapshot, RouterModule, Routes } from '@angular/router';
import { ZoneIndexComponent } from './zone-index/zone-index.component';
import { UserSettingsComponent } from './user-settings/user-settings.component';
import { ZoneViewComponent } from './zone-view/zone-view.component';

const routes: Routes = [
  { path: '', component: ZoneIndexComponent, title: 'Olympus' },
  {
    path: 'settings',
    component: UserSettingsComponent,
    title: 'Olympus: Settings',
  },
  {
    path: 'host/:host/zone/:zone',
    component: ZoneViewComponent,
    title: (route: ActivatedRouteSnapshot) => {
      return (
        'Olympus: ' +
        route.paramMap.get('host') +
        '.' +
        route.paramMap.get('zone')
      );
    },
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
