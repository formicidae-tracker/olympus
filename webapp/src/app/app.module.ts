import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule }    from '@angular/common/http';

import { AppComponent } from './app.component';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { RouterModule } from '@angular/router';
import { ROUTES } from './app.route';

import { HomeComponent } from './home/home.component';
import { ZoneComponent } from './zone/zone.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ZonePreviewComponent } from './zone-preview/zone-preview.component';

import { OlympusService } from '@services/olympus';
import { ClimateChartComponent } from './climate-chart/climate-chart.component';
import { AlarmComponent } from './alarm/alarm.component';
import { StateComponent } from './state/state.component';
import { VideoJsComponent } from './video-js/video-js.component';

@NgModule({
    imports: [
  		NgbModule,
		BrowserModule,
		HttpClientModule,
        RouterModule.forRoot(ROUTES, { relativeLinkResolution: 'legacy' })
	],
	declarations: [
		AppComponent,
		HomeComponent,
		ZoneComponent,
		PageNotFoundComponent,
		ZonePreviewComponent,
		ClimateChartComponent,
		AlarmComponent,
		StateComponent,
		VideoJsComponent
	],
	providers: [OlympusService],
	bootstrap: [AppComponent]
})
export class AppModule { }
