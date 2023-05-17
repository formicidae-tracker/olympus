import { Component, OnInit } from '@angular/core';
import { ThemeService } from '../core/services/theme.service';
import { Title } from '@angular/platform-browser';
import { ActivationEnd, ResolveEnd, Router } from '@angular/router';
import { map, filter } from 'rxjs';

@Component({
  selector: 'app-nav-toolbar',
  templateUrl: './nav-toolbar.component.html',
  styleUrls: ['./nav-toolbar.component.scss'],
})
export class NavToolbarComponent implements OnInit {
  public darkTheme: boolean = false;
  public title: string = 'Olympus';
  private _currentURL: string = '';

  constructor(private themeService: ThemeService, private router: Router) {}

  ngOnInit(): void {
    this.themeService.isDarkTheme.subscribe((dark) => (this.darkTheme = dark));
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
    this.themeService.setDarkTheme(dark);
  }

  public isRoot(): boolean {
    return this._currentURL == '/';
  }
}
