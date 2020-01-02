import React from 'react';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';
import Link from '@material-ui/core/Link';
import NavBar from './NavBar';
import UserTable from './UserTable'
import { API } from './Api';

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {' Copyright © '}
      {new Date().getFullYear()}{' '}
      <Link color="inherit" href={API.Url}>
        {API.Url}
      </Link>{' '}
    </Typography>
  );
}

export default function App() {
  return (
    <div>
      <NavBar>
      </NavBar>
      <Container>
        <UserTable />
      </Container>
      <Copyright />
    </div>
  );
}