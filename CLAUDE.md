# Door — Panel Frontend

## Stack
- Angular 17+ (standalone components)
- TypeScript (strict mode, sem any)
- SCSS para estilos
- HttpClient para requisições

## Padrões obrigatórios
- Sempre usar standalone components
- Nunca usar NgModules
- Organizar por feature (não por tipo de arquivo)
- Tipagem forte em tudo — interfaces para todos os modelos
- Sem `any` em nenhuma circunstância

## Design system
- Tema: dark
- Background principal: #0e0e10
- Superfície: #111114
- Inputs: fundo #0e0e10, borda 0.5px solid #252528
- Botão primário: fundo #f0f0f0, texto #0e0e10
- Border-radius: 4px inputs e botões, 6px containers
- Sem gradientes, sem sombras exageradas
- Labels em uppercase, font-size 11px, letter-spacing espaçado

## Config App (backend)
- Base URL dev: http://localhost:7171
- Auth: JWT (access token em memória, refresh token no localStorage)
- Endpoints:
    - POST /api/v1/auth/login
    - POST /api/v1/auth/register
    - POST /api/v1/auth/token
    - POST /api/v1/auth/logout

## Estrutura de pastas
src/app/
core/          ← services, interceptors, guards
features/      ← auth, dashboard, routes, logs, settings
shared/        ← components e pipes reutilizáveis

## Git
- Branch atual: feat/view
- Commitar ao final de cada etapa
- Mensagens no padrão: feat:, fix:, refactor: