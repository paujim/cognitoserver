let backendHost;
const hostname = window && window.location && window.location.hostname;

if (hostname === 'dev.realsite.com') {
  backendHost = 'https://dev.realsite.com';
} else if (hostname === 'staging.realsite.com') {
  backendHost = 'https://staging.realsite.com';
} else if (/^qa/.test(hostname)) {
  backendHost = `https://api.${hostname}`;
} else {
  backendHost = process.env.REACT_APP_BACKEND_HOST || 'http://localhost:5000';
}

const API_ROOT = `${backendHost}/api/`;
const headers = { Accept: "application/json" }

const fetchUsers = () => {
  return fetch(API_ROOT + "user/list", { headers })
}

export const USER_API = {
  Url: API_ROOT,
  ListUsers: fetchUsers,
}