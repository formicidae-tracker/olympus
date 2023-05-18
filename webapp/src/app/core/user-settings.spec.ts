import { cases } from 'jasmine-parameterized';

import { UserSettings } from 'src/app/core/user-settings';

describe('UserSettings', () => {
  it('should becreated with default values', () => {
    let e = new UserSettings();
    expect(e).toBeTruthy();
    expect(e.darkMode).toBe(false);
    expect(e.alarmSubscriptions).toEqual(new Set<string>([]));
  });

  it('should be created with non-default values', () => {
    let e = new UserSettings({
      darkMode: true,
      alarmSubscriptions: new Set<string>(['some.zone']),
    });
    expect(e).toBeTruthy();
    expect(e.darkMode).toEqual(true);
    expect(e.alarmSubscriptions).toEqual(new Set<string>(['some.zone']));
  });

  cases([
    [{}, '{"darkMode":false,"alarmSubscriptions":[]}'],
    [
      { darkMode: true, alarmSubscriptions: new Set<string>(['some.zone']) },
      '{"darkMode":true,"alarmSubscriptions":["some.zone"]}',
    ],
  ]).it('should serialize to JSON', ([args, jsondata]) => {
    expect(new UserSettings(args).serialize()).toEqual(jsondata);
  });

  cases([
    [{}, '{}'],
    [{}, '{"darkMode":false,"alarmSubscriptions":[]}'],
    [
      { darkMode: true, alarmSubscriptions: new Set<string>(['some.zone']) },
      '{"darkMode":true,"alarmSubscriptions":["some.zone"]}',
    ],
  ]).it('should deserialize from JSON', ([args, jsondata]) => {
    expect(UserSettings.deserialize(jsondata)).toEqual(new UserSettings(args));
  });

  it('should handle alarm subscription', () => {
    let s = new UserSettings();
    expect(s.hasSubscriptionToAlarmFrom('foo')).toBeFalse();
    expect(s.subscribeToAlarmFrom('foo')).toBeTrue();
    expect(s.hasSubscriptionToAlarmFrom('foo')).toBeTrue();
    expect(s.subscribeToAlarmFrom('foo')).toBeFalse();
    expect(s.hasSubscriptionToAlarmFrom('foo')).toBeTrue();
    expect(s.unsubscribeFromAlarmFrom('foo')).toBeTrue();
    expect(s.hasSubscriptionToAlarmFrom('foo')).toBeFalse();
    expect(s.unsubscribeFromAlarmFrom('foo')).toBeFalse();
  });
});
