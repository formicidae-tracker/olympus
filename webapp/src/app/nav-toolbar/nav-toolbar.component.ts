import { Component, OnInit } from '@angular/core';
import { ActivationEnd, ResolveEnd, Router } from '@angular/router';
import { map, filter } from 'rxjs';
import { UserSettingsService } from '../core/user-settings.service';

@Component({
  selector: 'app-nav-toolbar',
  templateUrl: './nav-toolbar.component.html',
  styleUrls: ['./nav-toolbar.component.scss'],
})
export class NavToolbarComponent implements OnInit {
  public darkTheme: boolean = false;
  public title: string = 'Olympus';
  private _currentURL: string = '';

  constructor(
    private settingService: UserSettingsService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.settingService
      .isDarkTheme()
      .subscribe((dark) => (this.darkTheme = dark));
    this.router.events
      .pipe(
        filter((e) => e instanceof ActivationEnd),
        map((e) => e as ActivationEnd)
      )
      .subscribe((event: ActivationEnd) => {
        this.title = (event.snapshot.title || 'Olympus').replace(
          /^Olympus: /,
          ''
        );
      });

    this.router.events
      .pipe(
        filter((e) => e instanceof ResolveEnd),
        map((e) => e as ResolveEnd)
      )
      .subscribe((event: ResolveEnd) => {
        this._currentURL = event.url || '/';
      });
  }

  public setDarkTheme(dark: boolean): void {
    this.settingService.setDarkTheme(dark);
  }

  public isRoot(): boolean {
    return this._currentURL == '/';
  }
}
