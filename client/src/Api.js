import Cookies from 'js-cookie'

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

const API_ROOT = `${backendHost}/api/`
const API_LIST_USERS = API_ROOT + "user/list"
const API_GET_TOKEN = API_ROOT + "token"

const fetchUsers = () => {
  let access_token = getAccessToken()
  return fetch(API_LIST_USERS, {
    method: "GET",
    headers: {
      "Authorization": "Bearer " + access_token,
      "Content-Type": "application/json",
    },
  })
}

const getToken = (username, password) => {
  let formData = new URLSearchParams();
  formData.append("username", username)
  formData.append("password", password)
  return fetch(API_GET_TOKEN, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: formData.toString(),
  })
}

const getAccessToken = () => {
  const access_token = Cookies.get('__access_token')
  return access_token
}

const setAccessToken = (access_token) => {
  Cookies.set('__access_token', access_token, { expires: 7 })
}

const removeAccessToken = () => {
  Cookies.remove('__access_token')
}

export const API = {
  Url: API_ROOT,
  ListUsersUrl: API_LIST_USERS,
  FetchListUsers: fetchUsers,
  GetTokenUrl: API_GET_TOKEN,
  FetchGetToken: getToken,
  GetAccessToken : getAccessToken,
  SetAccessToken: setAccessToken,
}

