import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ZoneIndexComponent } from './zone-index/zone-index.component';

const routes: Routes = [{ path: '', component: ZoneIndexComponent }];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
