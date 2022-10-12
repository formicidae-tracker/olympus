import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule }    from '@angular/common/http';

import { AppComponent } from './app.component';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { RouterModule } from '@angular/router';
import { ROUTES } from './app.routes';



import { HomeComponent } from './home/home.component';
import { ZoneComponent } from './zone/zone.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ZonePreviewComponent } from './zone-preview/zone-preview.component';

import { OlympusService } from '@services/olympus';
import { ClimateViewComponent } from './climate-view/climate-view.component';
import { StateComponent } from './state/state.component';
//import { VideoJsComponent } from './video-js/video-js.component';
import { AlarmListComponent } from './alarm-list/alarm-list.component';
import { ClimateChartComponent } from './climate-chart/climate-chart.component';
import { NgChartsModule } from 'ng2-charts';
import { LogsComponent } from './logs/logs.component';



@NgModule({
    imports: [
		RouterModule.forRoot(ROUTES),
		NgbModule,
		BrowserModule,
		HttpClientModule,
		NgChartsModule,
	],
	declarations: [
		AppComponent,
		HomeComponent,
		ZoneComponent,
		PageNotFoundComponent,
		ZonePreviewComponent,
		ClimateViewComponent,
		StateComponent,
//		VideoJsComponent,
		AlarmListComponent,
		ClimateChartComponent,
		LogsComponent,
	],
	providers: [OlympusService],
	bootstrap: [AppComponent]
})
export class AppModule { }
