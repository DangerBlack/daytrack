const YOUR_API_PORT = ":3000"

const api_url = 'http://10.0.2.2:3000/v1';

const user_url = `${api_url}/users`;
const magic_link_url = `${user_url}/magic-link`;

export { api_url, user_url, magic_link_url };