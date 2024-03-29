import { NgModule, isDevMode } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';

import { CoreModule } from './core/core.module';

import { NavToolbarComponent } from './nav-toolbar/nav-toolbar.component';
import { ZoneIndexComponent } from './zone-index/zone-index.component';
import { ZoneCardComponent } from './zone-index/zone-card/zone-card.component';
import { UserSettingsComponent } from './user-settings/user-settings.component';
import { ZoneViewComponent } from './zone-view/zone-view.component';
import { ClimateChartComponent } from './zone-view/climate-chart/climate-chart.component';
import { ClimateStatusComponent } from './zone-view/climate-status/climate-status.component';
import { TrackingPlayerComponent } from './zone-view/tracking-player/tracking-player.component';
import { AlarmsReportsComponent } from './zone-view/alarms-reports/alarms-reports.component';
import { TrackingStatusComponent } from './zone-view/tracking-status/tracking-status.component';

import { NgxEchartsModule } from 'ngx-echarts';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ServiceWorkerModule } from '@angular/service-worker';
import { LogIndexComponent } from './log-index/log-index.component';
import { VersionComponent } from './version/version.component';

@NgModule({
  declarations: [
    AppComponent,
    NavToolbarComponent,
    ZoneIndexComponent,
    ZoneCardComponent,
    UserSettingsComponent,
    ZoneViewComponent,
    ClimateChartComponent,
    ClimateStatusComponent,
    TrackingPlayerComponent,
    AlarmsReportsComponent,
    TrackingStatusComponent,
    LogIndexComponent,
    VersionComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    FormsModule,
    CoreModule,
    NgxEchartsModule.forRoot({
      echarts: () => import('echarts'),
    }),
    RouterModule,
    ServiceWorkerModule.register('ngsw-worker.js', {
      enabled: !isDevMode(),
      // Register the ServiceWorker as soon as the application is stable
      // or after 30 seconds (whichever comes first).
      registrationStrategy: 'registerWhenStable:30000',
    }),
  ],
  providers: [DatePipe],
  bootstrap: [AppComponent],
})
export class AppModule {}
