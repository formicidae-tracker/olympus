import { NotificationSettings } from '../core/notification-settings';

export class NotificationSettingUpdate {
  public endpoint: string = '';
  public settings: NotificationSettings = new NotificationSettings();

  public serialize(): string {
    let plain: any = Object.assign({}, this);
    plain.settings = this.settings.toPlain();
    return JSON.stringify(plain);
  }

  static fromPlain(plain: any): NotificationSettingUpdate {
    let res = new NotificationSettingUpdate();
    res.endpoint = plain.endpoint || '';
    res.settings = NotificationSettings.fromPlain(plain.settings || {});
    console.log(plain, res);
    return res;
  }

  static deserialize(jsonData: string): NotificationSettingUpdate {
    return NotificationSettingUpdate.fromPlain(JSON.parse(jsonData));
  }
}
