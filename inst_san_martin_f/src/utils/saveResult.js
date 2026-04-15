/** Normalizes API result objects from actions / handlers. */
export function isSaveSuccess(res) {
  if (res == null) return false;
  const c = Number(res.code);
  return c === 200 || c === 201;
}
