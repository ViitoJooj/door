import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, switchMap, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const authInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
  const authService = inject(AuthService);
  const token = authService.getAccessToken();

  const request = token ? cloneWithToken(req, token) : req;

  return next(request).pipe(
    catchError((error) => {
      // se voltou 401, tenta renovar o access token com o refresh token
      if (error instanceof HttpErrorResponse && error.status === 401) {
        return authService.refresh().pipe(
          switchMap((res) => next(cloneWithToken(req, res.access_token))),
          catchError((err) => {
            authService.logout();
            return throwError(() => err);
          })
        );
      }
      return throwError(() => error);
    })
  );
};

function cloneWithToken(req: HttpRequest<unknown>, token: string): HttpRequest<unknown> {
  return req.clone({
    setHeaders: { Authorization: `Bearer ${token}` }
  });
}
