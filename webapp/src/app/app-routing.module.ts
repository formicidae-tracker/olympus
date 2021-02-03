import { NgModule } from '@angular/core';
import { Routes,RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { LogsComponent } from './logs/logs.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ZoneComponent } from './zone/zone.component';

const routes: Routes = [
	{ path: '', component: HomeComponent },
    { path: 'host/:hostName/zone/:zoneName', component: ZoneComponent },
	{ path: 'logs', component: LogsComponent },
    { path: 'not-found', component: PageNotFoundComponent},
];

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
