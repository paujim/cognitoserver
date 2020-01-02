import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import InputAdornment from '@material-ui/core/InputAdornment';
import AccountCircle from '@material-ui/icons/AccountCircle';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';
import IconButton from '@material-ui/core/IconButton';

export default function FormDialog() {
    const [open, setOpen] = React.useState(true);
    const [showPassword, setShowPassword] = React.useState(false);
    const handleClose = () => {
        setOpen(false);
    };

    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
      };

    return (
        <Dialog
            maxWidth="xs"
            open={open} 
            onClose={handleClose}
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
                <Button onClick={handleClose} color="primary">
                    Get Token
                </Button>
            </DialogActions>
        </Dialog>
    );
}