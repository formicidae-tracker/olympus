import { Component, OnInit } from '@angular/core';
import { ThemeService } from '../core/services/theme.service';

@Component({
  selector: 'app-nav-toolbar',
  templateUrl: './nav-toolbar.component.html',
  styleUrls: ['./nav-toolbar.component.scss'],
})
export class NavToolbarComponent implements OnInit {
  public darkTheme: boolean = false;

  constructor(private themeService: ThemeService) {}

  ngOnInit(): void {
    this.themeService.isDarkTheme.subscribe((dark) => (this.darkTheme = dark));
  }

  public setDarkTheme(dark: boolean): void {
    this.themeService.setDarkTheme(dark);
  }
}
