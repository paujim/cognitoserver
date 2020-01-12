import React from 'react';
import { API } from '../utils/Api';
import UseFetch from './UseFetch'


export default function FetchUsers(dataHandler, errorHandler, isLoggedIn) {

    const onData = (data) => {
        if (data.error) {
            throw data.error
        }
        if (dataHandler !== undefined) {
            dataHandler(data.users)
        }
    }

    const onError = (error) => {
        if (errorHandler !== undefined) {
            errorHandler(error)
        }
    }

    let isPending = UseFetch(API.FetchListUsers, onData, onError, [isLoggedIn])

    return isPending;
}