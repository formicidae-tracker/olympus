import { NgModule } from '@angular/core';
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
import { ReportLogsComponent } from './zone-view/alarms-reports/report-logs/report-logs.component';
import { TrackingStatusComponent } from './zone-view/tracking-status/tracking-status.component';

import { NgxEchartsModule } from 'ngx-echarts';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

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
    ReportLogsComponent,
    TrackingStatusComponent,
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
  ],
  providers: [DatePipe],
  bootstrap: [AppComponent],
})
export class AppModule {}
