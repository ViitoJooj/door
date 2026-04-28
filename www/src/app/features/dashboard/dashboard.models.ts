export interface Application {
  ID: number;
  Url: string;
  Country: string;
  Created_by: number;
  Updated_at: string;
  Created_at: string;
}

export interface ApplicationOutput {
  success: boolean;
  message: string;
  data: {
    id: number;
    url: string;
    country: string;
    updated_at: string;
    created_at: string;
  };
}

export interface RequestLog {
  id: number;
  method: string;
  path: string;
  query_string: string;
  status_code: number;
  response_time_ms: number;
  ip: string;
  country: string;
  user_agent: string;
  referer: string;
  request_size: number;
  response_size: number;
  internal: boolean;
  created_at: string;
}

export interface RequestLogListResponse {
  success: boolean;
  message: string;
  data: RequestLog[];
}

export interface EnvVar {
  Id: number;
  Name: string;
  Value: string;
}

export interface EnvOutput {
  success: boolean;
  message: string;
  data: EnvVar;
}

export interface AdminUser {
  id: number;
  username: string;
  email: string;
  role: string;
  active: boolean;
  updated_at: string;
  created_at: string;
}

export interface UserListOutput {
  success: boolean;
  message: string;
  data: AdminUser[];
}

export interface AdminCreateUserOutput {
  success: boolean;
  message: string;
  data: {
    user: AdminUser;
    temporary_password: string;
  };
}

export interface UserOutput {
  success: boolean;
  message: string;
  data: AdminUser;
}

export interface RateLimitData {
  id: number;
  requests_per_second: number;
  burst: number;
  progressive_rate_limit: boolean;
  updated_at: string;
  created_at: string;
}

export interface RateLimitOutput {
  success: boolean;
  message: string;
  data: RateLimitData;
}

export interface ProtocolSettingsData {
  id: number;
  allowed_protocol: 'http' | 'https' | 'both';
  apply_scope: 'all' | 'external';
  updated_at: string;
  created_at: string;
}

export interface ProtocolSettingsOutput {
  success: boolean;
  message: string;
  data: ProtocolSettingsData;
}

export interface IPAccessData {
  id: number;
  ip: string;
  created_by: number;
  updated_by: number;
  updated_at: string;
  created_at: string;
}

export interface IPAccessListOutput {
  success: boolean;
  message: string;
  data: IPAccessData[];
}

export interface IPAccessOutput {
  success: boolean;
  message: string;
  data: IPAccessData;
}

export type SpecialRouteType = 'login' | 'register';

export interface SpecialRouteData {
  id: number;
  route_type: SpecialRouteType;
  path: string;
  max_distinct_requests: number;
  window_seconds: number;
  ban_seconds: number;
  enabled: boolean;
  created_by: number;
  updated_by: number;
  created_at: string;
  updated_at: string;
}

export interface SpecialRouteListOutput {
  success: boolean;
  message: string;
  data: SpecialRouteData[];
}

export interface SpecialRouteOutput {
  success: boolean;
  message: string;
  data: SpecialRouteData;
}

export type HealthStatus = 'healthy' | 'degraded' | 'unhealthy' | 'unknown';

export interface HealthOverviewData {
  status: HealthStatus;
  window_minutes: number;
  generated_at: string;
  total_requests: number;
  server_errors: number;
  client_errors: number;
  availability: number;
  server_error_rate: number;
  client_error_rate: number;
  average_latency_ms: number;
  p95_latency_ms: number;
  requests_per_minute: number;
  unique_ips: number;
  unique_paths: number;
}

export interface HealthOverviewOutput {
  success: boolean;
  message: string;
  data: HealthOverviewData;
}

export interface HealthRouteData {
  method: string;
  path: string;
  status: HealthStatus;
  window_minutes: number;
  last_seen_at: string;
  request_count: number;
  server_errors: number;
  client_errors: number;
  availability: number;
  server_error_rate: number;
  client_error_rate: number;
  average_latency_ms: number;
  p95_latency_ms: number;
  requests_per_minute: number;
}

export interface HealthRoutesOutput {
  success: boolean;
  message: string;
  data: HealthRouteData[];
}
