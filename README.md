# Instituto San Martín — Portal web

Aplicación full stack para la gestión académico-administrativa del instituto: usuarios con distintos roles, alumnos, docentes, carreras y catálogos auxiliares. El frontend consume una API REST en Go; la persistencia está en **MongoDB** (por ejemplo MongoDB Atlas).

---

## Qué hace el proyecto

- **Autenticación** mediante inicio de sesión; las rutas protegidas del panel exigen un **JWT** válido enviado en las cabeceras de las peticiones.
- **Autorización por rol** (por ejemplo administrador, administrativo, alumno): el menú y las pantallas visibles dependen del tipo de usuario.
- **CRUD y consultas** sobre carreras, alumnos, docentes, usuarios, roles, turnos, tipos de examen, modalidades de cursado y datos relacionados.
- **Reloj de referencia Argentina** (`America/Argentina/Buenos_Aires`): se muestra en el inicio y en la barra superior para usuarios con sesión activa; el backend puede persistir un documento de referencia en la colección `time` cuando se consulta el endpoint correspondiente.

La comunicación entre capas es **cliente → API JSON → MongoDB**. No se incluye en este documento ninguna credencial, contraseña ni cadena de conexión real: deben configurarse solo en archivos y entornos privados.

---

## Arquitectura del repositorio

| Carpeta | Rol |
|--------|-----|
| `inst_san_martin_b` | **Backend**: API en **Go** (Gorilla Mux), driver oficial de MongoDB, CORS habilitado. |
| `inst_san_martin_f` | **Frontend**: **React** (Create React App), React Router, Bootstrap, Axios. |

---

## Requisitos del entorno

| Componente | Versión / notas |
|------------|-----------------|
| **Go** | **1.18** o superior (según `go.mod`). |
| **Node.js** | Recomendado **18 LTS** (o 16+). Los scripts de arranque usan `openssl-legacy-provider` por compatibilidad con la versión actual de `react-scripts`. |
| **npm** | Incluido con Node (versión compatible con el lockfile del proyecto, si existe). |
| **MongoDB** | Cluster accesible por URI **mongodb+srv** (p. ej. Atlas). La base lógica usada por el código es la indicada en la variable `DB_NAME` (por defecto en código: **`san_martin`**). |

---

## Base de datos

### Nombre de la base

En la aplicación se utiliza la base **`san_martin`** (configurable mediante variables de entorno del backend).

### Colecciones utilizadas por el backend

| Colección | Uso resumido |
|-----------|----------------|
| `user` | Cuentas de acceso y datos vinculados al login. |
| `student` | Fichas de alumnos. |
| `teacher` | Fichas de docentes. |
| `degree` | Carreras / planes de estudio asociados. |
| `role` | Definición de roles del sistema. |
| `shift` | Turnos. |
| `attendance_mode` | Modalidad de cursado (asistencia / condición). |
| `test_type` | Tipo de examen. |
| `study_plan` | Planes de estudio. |
| `level` | Niveles académicos y orden jerárquico (`nivel`, `jerarquia`). |
| `modalidad` | Modalidades (dominio académico). |
| `titulo_habilitante` | Títulos habilitantes. |
| `time` | Documento singleton de hora de referencia Argentina (hora, fecha, día, mes, año); se actualiza cuando un cliente autenticado llama al endpoint de reloj. |

Los esquemas exactos de cada documento están definidos en los modelos y paquetes `db` del backend.

---

## Cómo funciona (flujo general)

1. El usuario abre el **frontend** (navegador).
2. Tras el **login**, el backend devuelve un **token JWT** que el cliente guarda (p. ej. en cookie) y reenvía como `Authorization: Bearer<token>` en las peticiones protegidas.
3. El **middleware** del backend valida el token y, según la ruta, comprueba la conexión a la base y los permisos del rol.
4. Las respuestas suelen seguir un formato JSON común (`message`, `code`, `ok`, `data`).
5. El endpoint **`GET /health`** permite comprobar que el servicio está en marcha (sin lógica de negocio pesada).

---

## Módulos y pantallas (frontend)

Rutas principales bajo `inst_san_martin_f/src` (muchas protegidas con el componente `Protected`):

- **Inicio público** (`/`) y **login** (`/login`).
- **Home** (`/home`): panel principal con reloj Argentina (si hay datos cargados).
- **Usuarios** (`/users`), **roles** (`/roles`), **blanqueo de contraseña** (`/passwordblank`), **cambio de contraseña** (`/passwordchange`).
- **Alumnos** (`/students`), **docentes** (`/teachers`), **carreras** (`/degree`).
- **Turnos** (`/shift`), **tipo de examen** (`/testtype`), **modalidad** (`/pursuetype`).
- **Cierre de sesión** (`/logout`).

La barra de navegación y el menú lateral muestran opciones según el **rol** decodificado del token.

---

## API backend (módulos de rutas)

Registrados en `inst_san_martin_b/routes/main.go` (orden aproximado): tiempo/reloj Argentina, usuarios, roles, turnos, modalidad de cursado, tipo de examen, carreras, jerarquía, planes de estudio, alumnos y docentes. Cada grupo tiene controladores en `controllers/`, lógica en `services/` y acceso a datos en `db/`.

---

## Configuración y ejecución del backend

1. Instalar **Go 1.18+**.
2. Desde la raíz del backend:

   ```bash
   cd inst_san_martin_b
   ```

3. Crear el archivo **`config/.env`** (no subir valores reales al repositorio). El programa lo carga desde `./config/.env` al iniciar. Debe definir al menos:

   | Variable | Descripción (sin valores sensibles) |
   |----------|----------------------------------------|
   | `DB_NAME` | Nombre de la base de datos en el cluster. |
   | `DB_USER` | Usuario de MongoDB con permisos sobre esa base. |
   | `DB_PASSWORD` | Contraseña de ese usuario. |
   | `DB_CLUSTER` | Host del cluster (sin el prefijo `mongodb+srv://`). |
   | `PORT` | Puerto HTTP del API (si se omite, en código se usa **8081**). |

4. Ejecutar:

   ```bash
   go run .
   ```

   Alternativa de compilación:

   ```bash
   go build -o san_martin_b .
   ./san_martin_b
   ```

5. Comprobar que el servidor responde, por ejemplo: `GET http://localhost:<PORT>/health` (ajustar host y puerto según tu `.env`).

---

## Configuración y ejecución del frontend

1. Instalar **Node.js** (recomendado 18 LTS).
2. En la carpeta del cliente:

   ```bash
   cd inst_san_martin_f
   npm install
   ```

3. Definir la URL base del API para Axios. Suele hacerse con un archivo **`.env`** en la raíz de `inst_san_martin_f` (no commitear secretos):

   ```env
   REACT_APP_BACKEND_URL=http://localhost:8081
   ```

   Ajustar el host y el puerto al mismo `PORT` que usa el backend.

4. Arranque en desarrollo:

   ```bash
   npm start
   ```

   La aplicación se sirve habitualmente en **http://localhost:3000**.

5. Build de producción:

   ```bash
   npm run build
   ```

   El resultado queda en la carpeta `build/`.

---

## Versiones principales de dependencias

### Backend (`inst_san_martin_b/go.mod`)

| Dependencia | Versión en el módulo |
|-------------|----------------------|
| Go | **1.18** |
| gorilla/mux | v1.8.0 |
| go.mongodb.org/mongo-driver | v1.9.1 |
| github.com/rs/cors | v1.8.2 |
| github.com/joho/godotenv | v1.4.0 |
| github.com/dgrijalva/jwt-go | v3.2.0+incompatible |
| golang.org/x/crypto | pin en go.mod (revisar `go.sum` para revisiones exactas) |

Para ver el árbol completo: `go list -m all` desde `inst_san_martin_b`.

### Frontend (`inst_san_martin_f/package.json`)

| Dependencia | Versión declarada |
|-------------|-------------------|
| react / react-dom | ^18.2.0 |
| react-router-dom | ^6.3.0 |
| react-scripts | ^2.1.3 |
| axios | ^0.27.2 |
| bootstrap | ^5.1.3 |
| react-bootstrap | ^2.4.0 |
| jwt-decode | ^3.1.2 |
| lodash | ^4.17.21 |
| @fortawesome/* | ^6.1.1 / ^0.1.18 |

Las versiones exactas instaladas pueden consultarse en `package-lock.json` si el repositorio lo incluye.

---

## Scripts npm del frontend

| Comando | Descripción |
|---------|-------------|
| `npm start` | Servidor de desarrollo (con `openssl-legacy-provider`). |
| `npm run build` | Genera la carpeta `build` para despliegue. |
| `npm test` | Ejecuta la suite de tests de CRA. |
| `npm run eject` | Expone la configuración interna de CRA (irreversible). |

---

## Seguridad y buenas prácticas

- No versionar **`config/.env`** del backend ni **`.env`** del frontend con datos reales.
- Rotar credenciales de base de datos y revisar permisos del usuario MongoDB en Atlas.
- En entornos productivos conviene externalizar secretos (JWT, cadena de conexión) a variables de entorno del servidor o a un gestor de secretos, en lugar de valores fijos en código.

---

## Licencia y autoría

Revisar el archivo `LICENSE` del repositorio si existe, y la política del instituto para uso y distribución del código.

Si algo de la configuración local difiere (puertos, nombre de base o cluster), ajustar solo los archivos `.env` y la documentación interna del equipo; este README describe el comportamiento previsto por el código actual.
