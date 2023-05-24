import { cases } from 'jasmine-parameterized';

import { UserSettings } from 'src/app/core/user-settings';

describe('UserSettings', () => {
  it('should becreated with default values', () => {
    let e = new UserSettings();
    expect(e).toBeTruthy();
    expect(e.darkMode).toBe(false);
    expect(e.subscriptions).toEqual(new Set<string>([]));
  });

  it('should be created with non-default values', () => {
    let e = new UserSettings({
      darkMode: true,
      subscriptions: new Set<string>(['some.zone']),
    });
    expect(e).toBeTruthy();
    expect(e.darkMode).toEqual(true);
    expect(e.subscriptions).toEqual(new Set<string>(['some.zone']));
  });

  cases([
    [
      {},
      '{"darkMode":false,"subscribeToAll":false,"notifyOnWarning":false,"subscriptions":[]}',
    ],
    [
      {
        darkMode: true,
        subscribeToAll: true,
        notifyOnWarning: true,
        subscriptions: new Set<string>(['some.zone']),
      },
      '{"darkMode":true,"subscribeToAll":true,"notifyOnWarning":true,"subscriptions":["some.zone"]}',
    ],
  ]).it('should serialize to JSON', ([args, jsondata]) => {
    expect(new UserSettings(args).serialize()).toEqual(jsondata);
  });

  cases([
    [{}, '{}'],
    [{}, '{"darkMode":false,"subscriptions":[]}'],
    [
      { darkMode: true, subscriptions: new Set<string>(['some.zone']) },
      '{"darkMode":true,"subscriptions":["some.zone"]}',
    ],
  ]).it('should deserialize from JSON', ([args, jsondata]) => {
    expect(UserSettings.deserialize(jsondata)).toEqual(new UserSettings(args));
  });

  it('should handle simple alarm subscription', () => {
    let s = new UserSettings();
    expect(s.hasSubscription('foo')).toBeFalse();
    expect(s.subscribeTo('foo')).toBeTrue();
    expect(s.hasSubscription('foo')).toBeTrue();
    expect(s.subscribeTo('foo')).toBeFalse();
    expect(s.hasSubscription('foo')).toBeTrue();
    expect(s.unsubscribeFrom('foo')).toBeTrue();
    expect(s.hasSubscription('foo')).toBeFalse();
    expect(s.unsubscribeFrom('foo')).toBeFalse();
  });

  it('should handle subscribeToAll', () => {
    let s = new UserSettings();
    expect(s.hasSubscription('foo')).toBeFalse();
    s.subscribeToAll = true;
    expect(s.hasSubscription('foo')).toBeTrue();
    expect(s.subscribeTo('foo')).toBeFalse();
    expect(s.unsubscribeFrom('foo')).toBeFalse();
    expect(s.subscriptions).toHaveSize(0);
    s.subscribeToAll = false;
    expect(s.hasSubscription('foo')).toBeFalse();
  });
});
