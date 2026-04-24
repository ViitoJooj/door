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
