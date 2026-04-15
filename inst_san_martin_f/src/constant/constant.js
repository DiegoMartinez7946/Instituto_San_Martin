export const ScreenPermission = {
  alumno: ["Main", "PasswordChange"],
  administrativo: ["Main", "Students", "Teachers", "PasswordChange", "Degree"],
  // Mismos módulos académicos que administrativo + gestión global (Usuarios, Roles, etc.)
  administrador: [
    "Main",
    "Roles",
    "Users",
    "Shift",
    "TestType",
    "PursueType",
    "PasswordBlank",
    "Students",
    "Teachers",
    "Degree"
  ],
};

export const LanguageTitles = {
  titles: ["Name", "Active", "Role", "Type", "Email"],
  titulos: ["Nombre", "Estado", "Rol", "Tipo", "Email"],
};
