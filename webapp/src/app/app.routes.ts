import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { LogsComponent } from './logs/logs.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ZoneComponent } from './zone/zone.component';

export const ROUTES: Routes = [
	{ path: '', component: HomeComponent },
	{ path: 'host/:hostName/zone/:zoneName', component: ZoneComponent },
	{ path: 'logs', component: LogsComponent },
	{ path: 'not-found', component: PageNotFoundComponent},
	{ path: '**', redirectTo: '/not-found'},
];
