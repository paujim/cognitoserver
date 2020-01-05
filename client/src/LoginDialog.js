import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import InputAdornment from '@material-ui/core/InputAdornment';
import AccountCircle from '@material-ui/icons/AccountCircle';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';
import IconButton from '@material-ui/core/IconButton';
import { API } from './Api';
import { LoginManager } from './LoginManager'

export default function FormDialog(props) {

    const [username, setUsername] = React.useState();
    const handleChangeUsername = event => {
        setUsername(event.target.value);
    };

    const [password, setPassword] = React.useState();
    const handleChangePassword = event => {
        setPassword(event.target.value);
    };

    const handleClose = () => {
        props.onClose();
    };
    const handleGetToken = () => {
        API.FetchGetToken(username, password)
            .then(resp => {
                return resp.json()
            })
            .then(data => {
                console.log(data)
                LoginManager.SetToken(data.access_token)
            })
            .catch(error => {
                console.log(error)
            })
            .finally(() => {
                props.onGetToken()
            })

    };

    const [showPassword, setShowPassword] = React.useState(false);
    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    return (
        <Dialog
            maxWidth="xs"
            open={props.open}
            // onClose={handleClose}
            disableEscapeKeyDown
            disableBackdropClick
            aria-labelledby="form-dialog-title">
            <DialogTitle id="form-dialog-title">Login</DialogTitle>
            <DialogContent >
                <TextField
                    // variant="outlined"
                    autoFocus
                    autoComplete="off"
                    id="username"
                    label="username"
                    type="text"
                    value={username}
                    onChange={handleChangeUsername}
                    fullWidth
                    InputProps={{
                        endAdornment: <InputAdornment position="end"><AccountCircle /></InputAdornment>,
                    }}
                />
                <TextField
                    autoComplete="off"
                    id="password"
                    label="password"
                    type={showPassword ? "text" : "password"}
                    value={password}
                    onChange={handleChangePassword}
                    fullWidth
                    InputProps={{
                        endAdornment: <InputAdornment position="end">
                            <IconButton
                                aria-label="toggle password visibility"
                                onClick={handleClickShowPassword}
                                edge="end"
                            >
                                {showPassword ? <Visibility /> : <VisibilityOff />}
                            </IconButton>
                        </InputAdornment>,
                    }}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} color="secondary">
                    Cancel
                </Button>
                <Button onClick={handleGetToken} color="primary">
                    Get Token
                </Button>
            </DialogActions>
        </Dialog>
    );
}