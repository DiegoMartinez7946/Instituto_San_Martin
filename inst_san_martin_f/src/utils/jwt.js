import jwt from 'jwt-decode';

/** Extrae el valor del JWT de document.cookie (solo el token, sin el prefijo `token=`). */
export const getAuthTokenFromCookie = () => {
  if (!document.cookie) return '';
  const parts = document.cookie.split(';');
  for (let i = 0; i < parts.length; i += 1) {
    const seg = parts[i].trim();
    if (seg.startsWith('token=')) {
      return seg.slice(6).trim();
    }
  }
  return '';
};

const safeDecode = (tokenString) => {
  if (!tokenString) return null;
  try {
    return jwt(tokenString);
  } catch {
    return null;
  }
};

/** Decodifica el JWT guardado en la cookie de sesión. */
export const decodeTokenFromCookie = () => safeDecode(getAuthTokenFromCookie());

/** Rol en mayúsculas para comparar con ADMINISTRADOR / ADMINISTRATIVO / ALUMNO (JWT puede variar mayúsculas). */
export const normalizeRoleFromToken = (decoded) => {
  if (!decoded || typeof decoded !== 'object') return '';
  const raw = decoded.type != null ? decoded.type : decoded.Type != null ? decoded.Type : decoded.role != null ? decoded.role : decoded.Role;
  return raw != null && String(raw).trim() !== '' ? String(raw).trim().toUpperCase() : '';
};

const decodeToken = (tokenString) => safeDecode(tokenString);

export default decodeToken;
