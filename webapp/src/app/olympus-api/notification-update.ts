import { NotificationSettings } from '../core/notification-settings';

export class NotificationSettingUpdate {
  public endpoint: string = '';
  public settings: NotificationSettings = new NotificationSettings();
}
