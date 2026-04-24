export interface UserData {
  id: number;
  username: string;
  email: string;
  role: string;
  active: boolean;
  updated_at: string;
  created_at: string;
}

export interface LoginRequest {
  username?: string;
  email?: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  data?: UserData;
}

export interface RegisterResponse {
  success: boolean;
  message: string;
  data?: UserData;
}
