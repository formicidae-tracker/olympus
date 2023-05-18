import { cases } from 'jasmine-parameterized';

import { TestBed } from '@angular/core/testing';

import { NotificationsSettingsService } from './notifications-settings.service';

describe('NotificationsSettingService', () => {
  let service: NotificationsSettingsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NotificationsSettingsService],
    });
    service = TestBed.inject(NotificationsSettingsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
