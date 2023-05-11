import { NgModule } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';

let UiComponents = [
  MatToolbarModule,
  MatButtonModule,
  MatIconModule,
  MatCardModule,
];

@NgModule({
  declarations: [],
  imports: [UiComponents],
  exports: [UiComponents],
})
export class CoreModule {}
