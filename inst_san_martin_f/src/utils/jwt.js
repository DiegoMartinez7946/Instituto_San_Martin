import jwt from 'jwt-decode';

const decodeToken = (token) => token ? jwt(token) : "";

export default decodeToken;

