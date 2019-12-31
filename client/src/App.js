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

const users = [
  createData('Frozenyoghurt', '2019-12-12 07:40:33', 'CONFIRMED', true),
  createData('Ssandwich', '2019-12-12 07:40:33', 'CONFIRMED', true),
  createData('Eclair', '2019-12-12 07:40:33', 'CONFIRMED', true),
  createData('Cupcake', '2019-12-12 07:40:33', 'CONFIRMED', false),
  createData('Gingerbread', '2019-12-12 07:40:33', 'UNCONFIRMED', true),
];
function createData(username, created, status, enabled) {
  return { username, created, status, enabled };
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