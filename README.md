# Sistema de Gestión de Condominios — API Backend

Proyecto Semestral · Aplicaciones Web II (TDI-601) · Período 2026-1
Universidad Laica Eloy Alfaro de Manabí

## ¿Qué es este producto?

API REST desarrollada en Go para la gestión integral de un condominio: cuotas y sanciones económicas a residentes (Obligaciones y Multas), control de vehículos y visitas (Parqueo), y gestión de áreas comunes y reservas (Área Social). Incluye autenticación JWT con roles (`admin` / `residente`) y persistencia mediante GORM sobre SQLite (desarrollo) o PostgreSQL (producción/Docker).

## Equipo — Paralelo A, Grupo B

| Integrante | Módulo |
|---|---|
| Winter Povea | Parqueo (vehículos, visitas, accesos) |
| Frank Montece | Área Social (áreas comunes, reservas, notificaciones) |
| Héctor Fernández | Obligaciones y Multas |

## Stack tecnológico

- **Lenguaje:** Go
- **Router:** [chi](https://github.com/go-chi/chi)
- **ORM:** GORM
- **Base de datos:** SQLite (local/desarrollo) · PostgreSQL (Docker/producción)
- **Autenticación:** JWT ([golang-jwt](https://github.com/golang-jwt/jwt)) + bcrypt
- **Testing:** paquete estándar `testing` + [testify](https://github.com/stretchr/testify) (`assert`, `require`)
- **Contenedores:** Docker (multi-stage build) + Docker Compose
- **CI/CD:** GitHub Actions (build → vet → test)

## Cómo correrlo

### Con Docker (recomendado)

```bash
docker compose up --build
```

Esto levanta la API junto con PostgreSQL y los seeders de datos iniciales, sin pasos manuales adicionales. Al arrancar por primera vez, la base se migra automáticamente y se siembran datos de ejemplo (usuarios, obligaciones, vehículos, visitas, accesos, áreas sociales). La API queda disponible en `http://localhost:8080`.

Para reiniciar desde cero (por ejemplo, tras cambiar el esquema):

```bash
docker compose down -v
docker compose up --build
```

### Localmente (sin Docker)

```bash
go mod download
go run cmd/api/main.go
```

Requiere un archivo `.env` en la raíz (ver `.env.example`) con al menos:

```
DB_DRIVER=sqlite
DB_DSN=
RUTA_DB=condominio.db
JWT_SECRETO=tu-secreto-aqui
PUERTO=:8080
```

### Verificar que levantó correctamente

```bash
curl http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}'
```

Debe responder `201 Created`.

### Usuarios de prueba (seed)

| Email | Password | Rol |
|---|---|---|
| `residente1@test.com` | `password123` | `admin` <!-- ⚠️ Winter: confirmar si el rol de este usuario es correcto, el email sugiere "residente" --> |
| `residente2@test.com` | `password123` | `residente` |
| `residente3@test.com` | `password123` | `residente` |

## Variables de entorno

Ver `.env.example` para la lista completa. Las más relevantes:

| Variable | Descripción |
|---|---|
| `DB_DRIVER` | `postgres` o `sqlite` (default) |
| `DB_DSN` | Cadena de conexión a PostgreSQL |
| `RUTA_DB` | Ruta del archivo SQLite (si `DB_DRIVER=sqlite`) |
| `JWT_SECRETO` | Clave de firma del JWT |
| `JWT_DURACION` | Duración del token (ej. `24h`) |
| `PUERTO` | Puerto HTTP (ej. `:8080`) |

## Arquitectura

El proyecto sigue una arquitectura en capas, compartida por los tres módulos:

```
Cliente (Postman / Frontend)
        │
        ▼
   ┌─────────┐
   │  Router │  chi + middleware (Logger, Recoverer, CORS)
   └────┬────┘
        ▼
   ┌─────────┐
   │  Auth   │  valida JWT, extrae usuarioID y rol al contexto
   └────┬────┘
        ▼
   ┌───────────┐
   │RequireRol │  autoriza según rol (solo en rutas de escritura)
   └────┬──────┘
        ▼
   ┌─────────┐
   │ Handler │  handlers_obligaciones_obligaciones.go, handlers_vehiculos_parqueo.go, handlers_area_social.go...
   └────┬────┘
        ▼
   ┌─────────┐
   │ Service │  lógica de negocio + validaciones (obligaciones.go, vehiculos_service.go, area_service.go...)
   └────┬────┘
        ▼
   ┌─────────┐
   │ Storage │  interfaz Repository → implementación GORM (obligaciones_sqlite.go, vehiculos_sqlite.go...)
   └────┬────┘
        ▼
   SQLite / PostgreSQL
```

Cada módulo define su propia interfaz de repositorio en `internal/storage/`, lo que permite probar los servicios con mocks/fakes sin depender de una base de datos real.

## Autenticación

Todos los endpoints bajo `/api/v1` (excepto `/auth/register` y `/auth/login`) requieren un JWT válido:

```
Authorization: Bearer {token}
```

| Endpoint | Método | Descripción |
|---|---|---|
| `/api/v1/auth/register` | POST | Registra un nuevo usuario (rol `residente` por defecto) |
| `/api/v1/auth/login` | POST | Autentica y devuelve un JWT con `usuario_id` y `rol` |

Existen dos roles: `residente` y `admin`. Las operaciones de lectura (GET) están abiertas a cualquier usuario autenticado; las operaciones de escritura (POST/PUT/DELETE) en los módulos de Obligaciones y Multas están restringidas a `admin`.

---

## Módulo Obligaciones y Multas — Héctor Fernández

Gestiona las cuotas de condominio (obligaciones mensuales y extraordinarias) y las sanciones económicas (multas) aplicadas a los residentes.

**Autenticación:** lectura abierta a cualquier autenticado; escritura restringida a rol `admin`.

### Obligaciones

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/obligaciones` | Cualquier autenticado | Lista todas las obligaciones registradas |
| GET    | `/api/v1/obligaciones/{id}` | Cualquier autenticado | Obtiene una obligación por su ID |
| POST   | `/api/v1/obligaciones` | admin | Crea una nueva obligación (mensual o extraordinaria) |
| PUT    | `/api/v1/obligaciones/{id}` | admin | Actualiza una obligación existente |
| DELETE | `/api/v1/obligaciones/{id}` | admin | Elimina una obligación |

**Ejemplo de body (POST/PUT):**
```json
{
  "residente_id": 1,
  "tipo": "mensual",
  "monto": 150.00,
  "periodo": "2026-07"
}
```

**Validaciones del service:** `residente_id` requerido (>0), `monto` mayor a 0, `periodo` no vacío (formato `YYYY-MM`), `tipo` debe ser `"mensual"` o `"extraordinaria"`. El campo `estado` lo asigna el sistema (`"pendiente"` al crear).

### Multas

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/multas` | Cualquier autenticado | Lista todas las multas registradas |
| GET    | `/api/v1/multas/{id}` | Cualquier autenticado | Obtiene una multa por su ID |
| POST   | `/api/v1/multas` | admin | Crea una nueva multa |
| PUT    | `/api/v1/multas/{id}` | admin | Actualiza una multa existente |
| DELETE | `/api/v1/multas/{id}` | admin | Elimina una multa |

**Ejemplo de body (POST/PUT):**
```json
{
  "residente_id": 1,
  "motivo": "Ruido excesivo en horario nocturno",
  "monto": 50.00
}
```

**Validaciones del service:** `motivo` requerido (no vacío).

**Relación GORM:** `Obligacion` y `Multa` tienen una relación Belongs-To hacia `Usuario` (`foreignKey:ResidenteID`).

---

## Módulo Parqueo — Winter Povea

Módulo de control de acceso y gestión vehicular. Cubre tres entidades relacionadas: Vehículos (registro de vehículos de residentes), Visitas (control de ingreso de vehículos visitantes, con código QR y registro de entrada/salida) y Accesos (bitácora de movimientos de vehículos registrados).

### Vehículos

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/vehiculos` | Cualquier autenticado | Lista todos los vehículos |
| GET    | `/api/v1/vehiculos/{id}` | Cualquier autenticado | Obtiene un vehículo por ID |
| POST   | `/api/v1/vehiculos` | admin | Crea un vehículo |
| PUT    | `/api/v1/vehiculos/{id}` | admin | Actualiza un vehículo |
| DELETE | `/api/v1/vehiculos/{id}` | admin | Elimina un vehículo |

**Ejemplo de body (POST/PUT):**
```json
{
  "residente_id": 1,
  "placa": "MKS-2841",
  "marca": "Nissan",
  "modelo": "Versa",
  "color": "Rojo"
}
```

### Visitas

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/visitas` | Cualquier autenticado | Lista todas las visitas |
| GET    | `/api/v1/visitas/{id}` | Cualquier autenticado | Obtiene una visita por ID |
| POST   | `/api/v1/visitas` | admin | Registra una nueva visita |
| PUT    | `/api/v1/visitas/{id}/entrada` | admin | Actualiza datos de entrada |
| PUT    | `/api/v1/visitas/{id}/salida` | admin | Registra la salida (expira el QR) |
| DELETE | `/api/v1/visitas/{id}` | admin | Elimina una visita |

**Ejemplo de body (POST/PUT):**
```json
{
  "condominio_id": 1,
  "residente_id": 1,
  "placa_visitante": "ABC123",
  "nombre_visitante": "Carlos Pérez",
  "motivo": "Visita familiar"
}
```

### Accesos

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/accesos` | Cualquier autenticado | Lista todos los accesos |
| GET    | `/api/v1/accesos/{id}` | Cualquier autenticado | Obtiene un acceso por ID |
| POST   | `/api/v1/accesos` | admin | Registra un movimiento (entrada/salida) |
| DELETE | `/api/v1/accesos/{id}` | admin | Elimina un registro de acceso |

**Ejemplo de body (POST/PUT):**
```json
{
  "condominio_id": 1,
  "residente_id": 1,
  "placa_visitante": "ABC123",
  "nombre_visitante": "Carlos Pérez",
  "motivo": "Visita familiar"
}
```

## Módulo Área Social — Frank Montece

Gestiona los espacios comunes del condominio (salones, canchas, piscina, etc.), las reservas que los residentes hacen sobre ellos, y las notificaciones internas asociadas (por ejemplo, confirmación de una reserva).

**Autenticación:** todas las rutas de este módulo requieren un JWT válido (`Authorization: Bearer {token}`). Actualmente no aplican restricciones por rol (`admin`/`residente`) a nivel de escritura — cualquier usuario autenticado puede crear, actualizar o eliminar.

### Áreas Sociales

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/areas-sociales` | Cualquier autenticado | Lista todas las áreas sociales |
| GET    | `/api/v1/areas-sociales/{id}` | Cualquier autenticado | Obtiene un área social por ID |
| POST   | `/api/v1/areas-sociales` | Cualquier autenticado | Crea una nueva área social |
| PUT    | `/api/v1/areas-sociales/{id}` | Cualquier autenticado | Actualiza un área social existente |
| DELETE | `/api/v1/areas-sociales/{id}` | Cualquier autenticado | Elimina un área social |

**Ejemplo de body (POST/PUT):**
```json
{
  "nombre": "Salón Principal",
  "descripcion": "Salón de eventos",
  "capacidad": 100,
  "activo": true
}
```

### Reservas de Áreas

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/reservas-areas` | Cualquier autenticado | Lista todas las reservas |
| GET    | `/api/v1/reservas-areas/{id}` | Cualquier autenticado | Obtiene una reserva por ID |
| POST   | `/api/v1/reservas-areas` | Cualquier autenticado | Crea una nueva reserva |
| PUT    | `/api/v1/reservas-areas/{id}` | Cualquier autenticado | Actualiza una reserva existente |
| DELETE | `/api/v1/reservas-areas/{id}` | Cualquier autenticado | Elimina una reserva |

**Ejemplo de body (POST/PUT):**
```json
{
  "area_id": 1,
  "residente_id": 1,
  "fecha_inicio": "2026-07-10T10:00:00Z",
  "fecha_fin": "2026-07-10T14:00:00Z",
  "proposito": "Cumpleaños",
  "numero_personas": 20
}
```

**Validaciones:** `area_id` y `residente_id` deben ser mayores a 0, `proposito` requerido, `numero_personas` mayor a 0. El campo `estado` lo asigna el sistema (`"pendiente"` al crear).

### Notificaciones

| Método | Endpoint | Rol requerido | Descripción |
|--------|----------|----------------|-------------|
| GET    | `/api/v1/notificaciones` | Cualquier autenticado | Lista todas las notificaciones |
| GET    | `/api/v1/notificaciones/{id}` | Cualquier autenticado | Obtiene una notificación por ID |
| POST   | `/api/v1/notificaciones` | Cualquier autenticado | Crea una nueva notificación |
| POST   | `/api/v1/notificaciones/{id}/leer` | Cualquier autenticado | Marca una notificación como leída |
| DELETE | `/api/v1/notificaciones/{id}` | Cualquier autenticado | Elimina una notificación |

**Ejemplo de body (POST):**
```json
{
  "residente_id": 1,
  "tipo": "reserva",
  "mensaje": "Tu reserva fue confirmada"
}
```

---

## Testing

```bash
go test ./... -cover
```

Cada integrante cubre las 3 capas (storage → service → handler) de una entidad representativa de su módulo, con mocks/fakes del repositorio y casos edge (error paths, inputs inválidos). Cobertura mínima exigida: 50% en la capa de service.
