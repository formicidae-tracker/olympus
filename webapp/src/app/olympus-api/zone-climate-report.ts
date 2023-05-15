import { ClimateState } from './climate-state';
import { Bounds } from './bounds';
import { Type, plainToClass } from 'class-transformer';

export class ZoneClimateReport {
  temperature?: number;
  humidity?: number;

  @Type(() => Bounds)
  temperature_bounds?: Bounds;
  @Type(() => Bounds)
  humidity_bounds?: Bounds;

  @Type(() => ClimateState)
  current?: ClimateState;
  @Type(() => ClimateState)
  current_end?: ClimateState;

  @Type(() => ClimateState)
  next?: ClimateState;
  @Type(() => ClimateState)
  next_end?: ClimateState;
  @Type(() => Date)
  next_time?: Date;

  static fromPlain(plain: any): ZoneClimateReport {
    return plainToClass(ZoneClimateReport, plain, {
      exposeDefaultValues: true,
    });
  }
}
