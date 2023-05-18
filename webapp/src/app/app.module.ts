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
import { ZoneStatusComponent } from './zone-view/zone-status/zone-status.component';
import { TrackingPlayerComponent } from './zone-view/tracking-player/tracking-player.component';
import { AlarmsReportsComponent } from './zone-view/alarms-reports/alarms-reports.component';

@NgModule({
  declarations: [
    AppComponent,
    NavToolbarComponent,
    ZoneIndexComponent,
    ZoneCardComponent,
    UserSettingsComponent,
    ZoneViewComponent,
    ClimateChartComponent,
    ZoneStatusComponent,
    TrackingPlayerComponent,
    AlarmsReportsComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    CoreModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
