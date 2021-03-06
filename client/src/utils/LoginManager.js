import Cookies from 'js-cookie'

const TOKEN_KEY = '__access_token'

const IsLoggedIn = () => {
    const access_token = Cookies.get(TOKEN_KEY)
    if (access_token === undefined || access_token === "") {
        return false
    }
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
    IsLoggedIn,
    LogOut,
    SetToken,
    GetToken,
}