import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatBadgeModule } from '@angular/material/badge';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatListModule } from '@angular/material/list';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatChipsModule } from '@angular/material/chips';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { VgCoreModule } from '@videogular/ngx-videogular/core';
import { VgControlsModule } from '@videogular/ngx-videogular/controls';
import { VgOverlayPlayModule } from '@videogular/ngx-videogular/overlay-play';
import { VgBufferingModule } from '@videogular/ngx-videogular/buffering';
import { VgStreamingModule } from '@videogular/ngx-videogular/streaming';

import { BoundedProgressBarComponent } from './bounded-progress-bar/bounded-progress-bar.component';
import { ZoneNotificationButtonComponent } from './zone-notification-button/zone-notification-button.component';

let UiComponents = [
  MatToolbarModule,
  MatButtonModule,
  MatIconModule,
  MatCardModule,
  MatProgressBarModule,
  MatBadgeModule,
  MatProgressSpinnerModule,
  MatExpansionModule,
  MatTableModule,
  MatPaginatorModule,
  MatButtonToggleModule,
  MatListModule,
  MatFormFieldModule,
  MatChipsModule,
  MatAutocompleteModule,
  MatSlideToggleModule,
];

let VideoComponents = [
  VgCoreModule,
  VgControlsModule,
  VgOverlayPlayModule,
  VgBufferingModule,
  VgStreamingModule,
];

@NgModule({
  declarations: [BoundedProgressBarComponent, ZoneNotificationButtonComponent],
  imports: [CommonModule, UiComponents, VideoComponents],
  exports: [
    UiComponents,
    VideoComponents,
    BoundedProgressBarComponent,
    ZoneNotificationButtonComponent,
  ],
})
export class CoreModule {}
