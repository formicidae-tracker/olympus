import { Component, OnInit } from '@angular/core';
import { TitleService } from '../core/services/title.service';

@Component({
  selector: 'app-user-settings',
  templateUrl: './user-settings.component.html',
  styleUrls: ['./user-settings.component.scss'],
})
export class UserSettingsComponent implements OnInit {
  constructor(private title: TitleService) {}

  ngOnInit() {
    this.title.setTitle('User Settings');
  }
}
