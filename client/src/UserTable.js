import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow }  from '@material-ui/core';
import Switch from '@material-ui/core/Switch';
import { IconButton, Button } from '@material-ui/core/';
import AddIcon from '@material-ui/icons/Add';
import { Card, CardContent, CardActions } from '@material-ui/core/';
import CircularProgress from '@material-ui/core/CircularProgress';
import { Dialog, DialogActions, DialogContent,DialogContentText, DialogTitle } from '@material-ui/core/';
import { USER_API } from './api-config';
import Snackbar from '@material-ui/core/Snackbar';

const useStyles = makeStyles(theme => ({
    circular: {
      position: 'absolute',
      right: '50%',
      display: 'flex',
      '& > * + *': {
        marginLeft: theme.spacing(2),
      },
    },
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

export default function UserTable() {
    const classes = useStyles();

    const [openDialog, setOpenDialog] = React.useState(false);
    const handleClickOpen = () => {
        setOpenDialog(true);
    };
    const handleCloseDialog = () => {
        setOpenDialog(false);
    };

    const [data, setData] = React.useState([]);
    const [isPending, setIsPending] = React.useState(false);
    const [isError, setIsError] = React.useState(false);
    const handleClose = (event, reason) => {
      if (reason === 'clickaway') {
        return;
      }
      setIsError(false);
    };

    React.useEffect(() => {
        const fetchData = async () => {
          const response = await USER_API.ListUsers()
          console.log(response)
          let data = await response.json()
          console.log(data)
          setIsPending(false)
          if (!response.ok){
            data=[]
            setIsError(true)
          }
          setData (data);
        };
        setIsPending (true)
        fetchData().catch(error => {
          console.log(error)
          setIsPending(false)
          setIsError(true)
          setData([]);
        })
      }, []);
    
    let progress =  isPending ? <div className={classes.circular}><CircularProgress /></div> : <div></div>;

    return (
    <Card>
        <CardContent>
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
                        {progress}
                        {data.map(row => (
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
        </CardContent>
        <CardActions>
            <IconButton aria-label="add user" onClick={handleClickOpen}>
                <AddIcon/>
            </IconButton>
        </CardActions>
        <Dialog open={openDialog} onClose={handleCloseDialog} aria-labelledby="form-dialog-title">
            <DialogTitle id="form-dialog-title">Add User</DialogTitle>
            <DialogContent>
            <DialogContentText>
                To add an user, please enter a username and password here.
            </DialogContentText>

            </DialogContent>
            <DialogActions>
          <Button onClick={handleCloseDialog} color="primary">
            Cancel
          </Button>
          <Button onClick={handleCloseDialog} color="primary">
            Add
          </Button>
        </DialogActions>
        </Dialog>
        <Snackbar 
            anchorOrigin={{vertical: 'bottom',horizontal: 'left'}}
            open={isError}
            autoHideDuration={6000}
            onClose={handleClose}
            ContentProps={{'aria-describedby': 'message-id'}}
            message={<span id="message-id">Error Fetching data</span>}
            >
        </Snackbar> 
    </Card>

  );
}