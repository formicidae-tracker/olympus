import { NotificationSettings } from '../core/notification-settings';

export class NotificationSettingsUpdate {
  constructor(
    public endpoint: string = '',
    public settings: NotificationSettings = new NotificationSettings()
  ) {}

  public serialize(): string {
    let plain: any = Object.assign({}, this);
    plain.settings = this.settings.toPlain();
    return JSON.stringify(plain);
  }

  static fromPlain(plain: any): NotificationSettingsUpdate {
    let res = new NotificationSettingsUpdate();
    res.endpoint = plain.endpoint || '';
    res.settings = NotificationSettings.fromPlain(plain.settings || {});
    console.log(plain, res);
    return res;
  }

  static deserialize(jsonData: string): NotificationSettingsUpdate {
    return NotificationSettingsUpdate.fromPlain(JSON.parse(jsonData));
  }
}
