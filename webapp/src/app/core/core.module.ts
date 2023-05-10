import { NgModule } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

let UiComponents = [MatToolbarModule, MatButtonModule, MatIconModule];

@NgModule({
  declarations: [],
  imports: [UiComponents],
  exports: [UiComponents],
})
export class CoreModule {}
