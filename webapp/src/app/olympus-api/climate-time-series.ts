import { Type, plainToClass } from 'class-transformer';

export class Point {
  constructor(public x: number = 0.0, public y: number = 0.0) {}
}

export class ClimateTimeSeries {
  units: string = '';
  @Type(() => Date)
  reference: Date = new Date(0);
  @Type(() => Point)
  humidity: Point[] = [];
  @Type(() => Point)
  temperature: Point[] = [];
  @Type(() => Point)
  temperatureAux: Point[][] = [];

  static fromPlain(plain: any): ClimateTimeSeries {
    return plainToClass(ClimateTimeSeries, plain, {
      exposeDefaultValues: true,
    });
  }
}
