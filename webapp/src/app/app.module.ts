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

@NgModule({
  declarations: [
    AppComponent,
    NavToolbarComponent,
    ZoneIndexComponent,
    ZoneCardComponent,
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
