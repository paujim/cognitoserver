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
      {' Copyright © '}
      {new Date().getFullYear()}{' '}
      <Link color="inherit" href={USER_API.Url}>
        {USER_API.Url}
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