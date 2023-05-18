import { Component, OnInit } from '@angular/core';
import { UserSettingsService } from './core/user-settings.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  public darkTheme: boolean = false;
  constructor(private settings: UserSettingsService) {}

  ngOnInit(): void {
    this.settings.isDarkTheme().subscribe((dark) => (this.darkTheme = dark));
  }
}
