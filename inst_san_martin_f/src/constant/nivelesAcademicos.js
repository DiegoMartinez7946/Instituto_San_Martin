/** Valores por defecto si GET /level aún no cargó (colección Mongo `level`) */
export const ORDEN_POR_NIVEL_DEFAULT = {
  inicial: 1,
  primario: 2,
  secundario: 3,
  terciario: 4,
  universitario: 4
};

export const NIVELES_ORDENADOS = Object.keys(ORDEN_POR_NIVEL_DEFAULT);

/** Mapa nivel → orden desde filas del API (`nivel`, `jerarquia`) */
export const mapaOrdenDesdeLevels = (levels) => {
  if (!Array.isArray(levels) || !levels.length) return null;
  const m = {};
  for (const row of levels) {
    const n = row.nivel != null ? String(row.nivel).toLowerCase().trim() : '';
    if (n) m[n] = Number(row.jerarquia);
  }
  return Object.keys(m).length ? m : null;
};

/**
 * Orden jerárquico del nivel (mayor = más avanzado).
 * @param {string} nivel
 * @param {Record<string, number>|null|undefined} mapaOpcional desde `globalState.levels`
 */
export const jerarquiaDeNivel = (nivel, mapaOpcional) => {
  const mapa =
    mapaOpcional && typeof mapaOpcional === 'object' && Object.keys(mapaOpcional).length
      ? mapaOpcional
      : ORDEN_POR_NIVEL_DEFAULT;
  const k = nivel != null ? String(nivel).toLowerCase().trim() : '';
  if (!k) return null;
  return Object.prototype.hasOwnProperty.call(mapa, k) ? mapa[k] : null;
};

export const etiquetaNivel = (nivel) => {
  if (!nivel) return '';
  return nivel.charAt(0).toUpperCase() + nivel.slice(1).toLowerCase();
};
