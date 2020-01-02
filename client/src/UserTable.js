import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@material-ui/core';
import Switch from '@material-ui/core/Switch';

const useStyles = makeStyles(theme => ({
  table: {
    minWidth: 650,
  },
}));

const mock_users = [
  createData('Frozenyoghurt', '2019-12-12 07:40:33', 'CONFIRMED', true),
  createData('Ssandwich', '2019-12-12 07:40:33', 'CONFIRMED', true),
  createData('Eclair', '2019-12-12 07:40:33', 'CONFIRMED', true),
  createData('Cupcake', '2019-12-12 07:40:33', 'CONFIRMED', false),
  createData('Gingerbread', '2019-12-12 07:40:33', 'UNCONFIRMED', true),
];
function createData(username, created, status, enabled) {
  return { username, created, status, enabled };
}

const getUsers = (users) => {
  if (users === undefined) {
    return []
  }
  return users
}

export default function UserTable(props) {
  const classes = useStyles();
  return (
    <TableContainer>
      <Table className={classes.table} aria-label="user table">
        <TableHead>
          <TableRow>
            <TableCell>Username</TableCell>
            <TableCell>Created</TableCell>
            <TableCell>Status</TableCell>
            <TableCell>Enabled</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {getUsers(props.users).map(row => (
            <TableRow key={row.username}>
              <TableCell component="th" scope="row">{row.username}</TableCell>
              <TableCell>{row.created}</TableCell>
              <TableCell>{row.status}</TableCell>
              <TableCell>
                <Switch checked={row.enabled} color="primary" value="checkedEnabled" inputProps={{ 'aria-label': 'primary checkbox' }} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}