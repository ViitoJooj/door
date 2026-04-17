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
