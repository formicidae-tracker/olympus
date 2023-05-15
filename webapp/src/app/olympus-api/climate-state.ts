import { plainToClass } from 'class-transformer';

export class ClimateState {
  public name: string = '';
  public temperature?: number;
  public humidity?: number;
  public wind?: number;
  public visible_light?: number;
  public uv_light?: number;

  static fromPlain(plain: any): ClimateState {
    return plainToClass(ClimateState, plain, {
      exposeDefaultValues: true,
    });
  }
}
