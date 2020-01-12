import React from 'react';
import Typography from '@material-ui/core/Typography';
import Link from '@material-ui/core/Link';
import NavBar from './components/NavBar';
import { API }   from './utils/Api';

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
      <NavBar />
      <Copyright />
    </div>
  );
}