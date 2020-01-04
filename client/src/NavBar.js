import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import LoginDialog from './LoginDialog'
import SecurityRoundedIcon from '@material-ui/icons/SecurityRounded';
import {LoginManager} from './LoginManager'


const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
  close: {
    padding: theme.spacing(0.5),
  },
}));

export default function NavBar(props) {
  const classes = useStyles();
  const [showLogin, setShowLogin] = React.useState(false);

  const handleLogin = () => {
    if (LoginManager.IsLogin()) {
      LoginManager.LogOut()
    }
    else {
      setShowLogin(true);
    }
  };



  let login = <LoginDialog />

  return (
    <div className={classes.root}>
      <AppBar position="static">
        <Toolbar>
          <SecurityRoundedIcon />
          <Typography variant="h6" className={classes.title}>
            Admin
          </Typography>
          <Button color="inherit" onClick={handleLogin} >{LoginManager.IsLogin() ? "Logout" : "Login"}</Button>
        </Toolbar>
        {showLogin ? login : null}
      </AppBar>
    </div>
  );
}
