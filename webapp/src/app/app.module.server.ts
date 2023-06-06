import { NgModule } from '@angular/core';
import { ServerModule } from '@angular/platform-server';

import { AppModule } from './app.module';
import { AppComponent } from './app.component';
import { Routes, RouterModule } from '@angular/router';

import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { AppShellComponent } from './app-shell/app-shell.component';

import {
  LocalStorageService,
  NullLocalStorageService,
} from './core/local-storage.service';
import {
  NetworkStatusService,
  ServerNetworkStatusService,
} from './core/network-status.service';

const routes: Routes = [{ path: 'shell', component: AppShellComponent }];

@NgModule({
  imports: [
    AppModule,
    ServerModule,
    RouterModule.forRoot(routes),
    MatProgressSpinnerModule,
  ],
  bootstrap: [AppComponent],
  declarations: [AppShellComponent],
  providers: [
    { provide: LocalStorageService, useClass: NullLocalStorageService },
    { provide: NetworkStatusService, useClass: ServerNetworkStatusService },
  ],
})
export class AppServerModule {}
