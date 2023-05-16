import { NgModule } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { BoundedProgressBarComponent } from './bounded-progress-bar/bounded-progress-bar.component';

let UiComponents = [
  MatToolbarModule,
  MatButtonModule,
  MatIconModule,
  MatCardModule,
  MatDividerModule,
  MatProgressBarModule,
];

@NgModule({
  declarations: [BoundedProgressBarComponent],
  imports: [UiComponents],
  exports: [UiComponents, BoundedProgressBarComponent],
})
export class CoreModule {}
