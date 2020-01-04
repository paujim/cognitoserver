import {LoginManager} from './LoginManager'

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

const FetchListUsers = () => {
  let access_token = LoginManager.GetToken()
  return fetch(API_LIST_USERS, {
    method: "GET",
    headers: {
      "Authorization": "Bearer " + access_token,
      "Content-Type": "application/json",
    },
  })
}

const FetchGetToken = (username, password) => {
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


export const API = {
  Url: API_ROOT,
  ListUsersUrl: API_LIST_USERS,
  FetchListUsers,
  GetTokenUrl: API_GET_TOKEN,
  FetchGetToken,
}

