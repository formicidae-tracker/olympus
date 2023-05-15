export class ClimateState {
  public name: string = '';
  public temperature?: number;
  public humidity?: number;
  public wind?: number;
  public visible_light?: number;
  public uv_light?: number;

  static fromPlain(plain: any): ClimateState | undefined {
    if (plain == undefined) {
      return undefined;
    }
    let res = new ClimateState();
    res.name = plain.name || '';
    res.temperature = plain.temperature;
    res.humidity = plain.humidity;
    res.wind = plain.wind;
    res.visible_light = plain.visible_light;
    res.uv_light = plain.uv_light;
    return res;
  }
}
