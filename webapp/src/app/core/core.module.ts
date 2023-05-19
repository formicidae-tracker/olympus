import { NgModule } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatBadgeModule } from '@angular/material/badge';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { VgCoreModule } from '@videogular/ngx-videogular/core';
import { VgControlsModule } from '@videogular/ngx-videogular/controls';
import { VgOverlayPlayModule } from '@videogular/ngx-videogular/overlay-play';
import { VgBufferingModule } from '@videogular/ngx-videogular/buffering';
import { VgStreamingModule } from '@videogular/ngx-videogular/streaming';

import { BoundedProgressBarComponent } from './bounded-progress-bar/bounded-progress-bar.component';

let UiComponents = [
  MatToolbarModule,
  MatButtonModule,
  MatIconModule,
  MatCardModule,
  MatProgressBarModule,
  MatBadgeModule,
  MatProgressSpinnerModule,
];

let VideoComponents = [
  VgCoreModule,
  VgControlsModule,
  VgOverlayPlayModule,
  VgBufferingModule,
  VgStreamingModule,
];

@NgModule({
  declarations: [BoundedProgressBarComponent],
  imports: [UiComponents, VideoComponents],
  exports: [UiComponents, VideoComponents, BoundedProgressBarComponent],
})
export class CoreModule {}
