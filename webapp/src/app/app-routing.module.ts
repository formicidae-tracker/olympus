import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ZoneIndexComponent } from './zone-index/zone-index.component';
import { UserSettingsComponent } from './user-settings/user-settings.component';
import { ZoneViewComponent } from './zone-view/zone-view.component';

const routes: Routes = [
  { path: '', component: ZoneIndexComponent },
  { path: 'settings', component: UserSettingsComponent },
  { path: 'host/:host/zone/:zone', component: ZoneViewComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
