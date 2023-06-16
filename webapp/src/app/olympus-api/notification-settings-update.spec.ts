import { cases } from 'jasmine-parameterized';

import { NotificationSettings } from '../core/notification-settings';
import { NotificationSettingUpdate } from './notification-settings-update';

import testData from './unit-testdata/NotificationSettingsUpdate.json';

describe('NotificationSettingsUpdate', () => {
  it('should create', () => {
    expect(new NotificationSettingUpdate()).toBeTruthy();
  });

  it('should parse and serialize to proper JSON', () => {
    for (const plain of testData) {
      let u: NotificationSettingUpdate =
        NotificationSettingUpdate.fromPlain(plain);
      expect(u.endpoint).toEqual(plain.endpoint || '');
      if (plain.settings != undefined) {
        expect(u.settings).toEqual(
          NotificationSettings.fromPlain(plain.settings)
        );
      }
    }
  });

  it('should serialize correctly to JSON', () => {
    for (const plain of testData) {
      if (plain.endpoint == undefined) {
        continue;
      }

      let u = NotificationSettingUpdate.fromPlain(plain);

      expect(JSON.parse(u.serialize())).toEqual(plain);
    }
  });
});
