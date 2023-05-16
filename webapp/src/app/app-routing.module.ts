import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ZoneIndexComponent } from './zone-index/zone-index.component';
import { UserSettingsComponent } from './user-settings/user-settings.component';

const routes: Routes = [
  { path: '', component: ZoneIndexComponent },
  { path: 'settings', component: UserSettingsComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
