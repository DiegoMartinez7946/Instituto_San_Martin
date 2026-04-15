import clientAxios from "../../../config/axios";

export const getArgentinaPortalTime = async (dispatch) => {
  const access_token = document.cookie.replace("token=", "");
  if (!access_token) {
    return;
  }
  const options = {
    headers: {
      Authorization: `Bearer${access_token}`,
    },
  };
  try {
    const result = await clientAxios.get("/time/argentina", options);
    const row = result.data && result.data.data;
    if (row) {
      dispatch({
        type: "SET_PORTAL_TIME",
        payload: row,
      });
    }
  } catch {
    /* token inválido o red: no actualizar */
  }
};
