const YOUR_API_PORT = ":3000"

const api_url = 'http://10.0.2.2:3000/v1';

const user_url = `${api_url}/users`;
const user_self_url = `${user_url}/self`;
const magic_link_url = `${user_url}/magic-link`;

const track_url = `${api_url}/tracks`;
const event_url = `${api_url}/events`

export { api_url, user_url, magic_link_url, track_url, event_url,user_self_url };