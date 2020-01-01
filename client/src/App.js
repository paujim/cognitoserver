import React from 'react';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';
import Link from '@material-ui/core/Link';
import NavBar from './NavBar';
import UserTable from './UserTable'
import { USER_API } from './api-config';

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {USER_API.Url}{' Copyright © '}
      <Link color="inherit" href="https://material-ui.com/">
        Your Website
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
      
    </Typography>
  );
}

export default function App() {
  return (
  <div> 
    <NavBar>
    </NavBar> 
    <Container>
      <UserTable/>
    </Container>
    <Copyright/>
  </div>
  );
}