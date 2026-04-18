export const ScreenPermission = {
  alumno: ["Main", "PasswordChange"],
  docente: ["Main", "PasswordChange"],
  administrativo: [
    "Main",
    "Students",
    "Teachers",
    "PasswordChange",
    "Degree",
    "StudyPlan",
    "Shift",
    "TestType",
    "PursueType",
  ],
  /** Solo cuentas de sistema (user), roles globales y blanqueo; sin carreras/turnos/etc. */
  administrador: ["Main", "Roles", "Users", "PasswordBlank", "PasswordChange"],
};

export const LanguageTitles = {
  titles: ["Name", "Active", "Role", "Type", "Email"],
  titulos: ["Nombre", "Estado", "Rol", "Tipo", "Email"],
};
