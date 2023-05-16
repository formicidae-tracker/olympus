import { Component, OnInit } from '@angular/core';
import { ThemeService } from './core/services/theme.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  public darkTheme: boolean = false;
  constructor(private theme: ThemeService) {}

  ngOnInit(): void {
    this.theme.isDarkTheme.subscribe((dark) => (this.darkTheme = dark));
  }
}
