import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { NodeIndexComponent } from './node-index/node-index.component';

const routes: Routes = [{ path: '', component: NodeIndexComponent }];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
