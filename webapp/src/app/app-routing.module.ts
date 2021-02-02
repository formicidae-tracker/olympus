import { NgModule } from '@angular/core';
import { Routes,RouterModule } from '@angular/router';
import { environment } from '@environments/environment';
import { AlarmListComponent } from './alarm-list/alarm-list.component';
import { HomeComponent } from './home/home.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ZoneComponent } from './zone/zone.component';

const routes: Routes = [
	{ path: '', component: HomeComponent },
    { path: 'host/:hostName/zone/:zoneName', component: ZoneComponent },
    { path: 'not-found', component: PageNotFoundComponent},
];

if (environment.production == false ) {
	routes.push({ path: 'debug-alarm-list', component: AlarmListComponent });
}

routes.push({ path: '**', redirectTo: '/not-found'});


@NgModule({
  declarations: [],
	imports: [
		RouterModule.forRoot(routes),
	],
	exports: [
		RouterModule,
	]
})
export class AppRoutingModule { }
