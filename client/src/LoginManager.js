import Cookies from 'js-cookie'

const TOKEN_KEY = '__access_token'

const IsLogin = () => { 
    const access_token = Cookies.get(TOKEN_KEY)
    console.log("access_token:")
    console.log(access_token)
    if (access_token === undefined || access_token === "")
        return false
    return true 
}

const LogOut = () => { 
    Cookies.remove(TOKEN_KEY)
}

const SetToken = (token) => {
    Cookies.set(TOKEN_KEY, token, { expires: 7 })
}

const GetToken = () => {
    let access_token = Cookies.get(TOKEN_KEY)
    return access_token
  }

export const LoginManager = {
    IsLogin,
    LogOut,
    SetToken,
    GetToken,
  }