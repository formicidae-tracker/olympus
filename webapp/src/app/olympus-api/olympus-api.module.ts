import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

@NgModule({
  declarations: [],
  imports: [CommonModule],
})
export class OlympusApiModule {}

export interface Bounds {
  min?: number;
  max?: number;
}

export interface ClimateState {
  name?: string;
  temperature?: number;
  humidity?: number;
  wind?: number;
  visible_light?: number;
  uv_light?: number;
}

export interface ZoneClimateReport {
  temperature?: number;
  humidity?: number;
  temperature_bounds: Bounds;
  humidity_bounds: Bounds;
  active_warnings: number;
  active_emergencies: number;
  current?: ClimateState;
  current_end?: ClimateState;
  next?: ClimateState;
  next_end?: ClimateState;
  next_time?: string;
}

export interface StreamInfo {
  experiment_name: string;
  stream_URL: string;
  thumbnail_URL: string;
}

export interface ZoneReportSummary {
  host: string;
  name: string;
  climate?: ZoneClimateReport;
  stream?: StreamInfo;
}
