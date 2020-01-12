import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { IconButton, Button } from '@material-ui/core/';
import AddIcon from '@material-ui/icons/Add';
import CloseIcon from '@material-ui/icons/Close';
import { Card, CardContent, CardActions } from '@material-ui/core/';
import Snackbar from '@material-ui/core/Snackbar';
import CircularProgress from '@material-ui/core/CircularProgress';
import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@material-ui/core/';
import { API } from '../utils/Api';
import UserTable from './UserTable';
import UseFetch from '../utils/UseFetch'

const useStyles = makeStyles(theme => ({
  circular: {
    position: 'absolute',
    right: '50%',
    top: '155px',
    display: 'flex',
    '& > * + *': {
      marginLeft: theme.spacing(2),
    },
  },
}));

export default function FetchUserTable() {
  const classes = useStyles();

  const [openDialog, setOpenDialog] = React.useState(false);
  const handleClickOpen = () => {
    setOpenDialog(true);
  };
  const handleCloseDialog = () => {
    setOpenDialog(false);
  };

  const dataHandler = (data) => {
    console.log(data)

    if (data.error){
      throw data.error
    }

    setData(data.users)
    setHasFailed(false)
  }

  const errorHandler = (error) => {
    console.log(error)
    setData([])
    setHasFailed(true)
  }

  const [data, setData] = React.useState([]);
  const isPending = UseFetch(API.FetchListUsers, dataHandler, errorHandler)
  const [hasFailed, setHasFailed] = React.useState(false);
  const handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setHasFailed(false);
  };

  let progress = isPending ? <div className={classes.circular}><CircularProgress /></div> : null;

  return (
    <Card>
      <CardContent>
        {progress}
        <UserTable users={data} />
      </CardContent>
      <CardActions>
        <IconButton aria-label="add user" onClick={handleClickOpen}>
          <AddIcon />
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
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        open={hasFailed}
        autoHideDuration={6000}
        onClose={handleClose}
        ContentProps={{ 'aria-describedby': 'message-id' }}
        message={<span id="message-id">Error Fetching data</span>}
        action={[
          <IconButton
            key="close"
            aria-label="close"
            color="inherit"
            className={classes.close}
            onClick={handleClose}
          >
            <CloseIcon />
          </IconButton>,
        ]}
      >
      </Snackbar>
    </Card>

  );
}