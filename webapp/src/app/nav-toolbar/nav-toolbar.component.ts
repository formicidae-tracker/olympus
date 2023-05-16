import { Component, OnInit } from '@angular/core';
import { ThemeService } from '../core/services/theme.service';
import { TitleService } from '../core/services/title.service';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-nav-toolbar',
  templateUrl: './nav-toolbar.component.html',
  styleUrls: ['./nav-toolbar.component.scss'],
})
export class NavToolbarComponent implements OnInit {
  public darkTheme: boolean = false;
  public title: string = 'Olympus';

  constructor(
    private themeService: ThemeService,
    private titleService: TitleService,
    private pageTitle: Title
  ) {}

  ngOnInit(): void {
    this.themeService.isDarkTheme.subscribe((dark) => (this.darkTheme = dark));
    this.titleService.title.subscribe((title) => {
      this.pageTitle.setTitle(title);
      this.title = title;
    });
  }

  public setDarkTheme(dark: boolean): void {
    this.themeService.setDarkTheme(dark);
  }
}
