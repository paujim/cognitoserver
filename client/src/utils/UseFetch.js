import React from 'react';

export default function UseFetch(fetchFunction, dataHandler, errorHandler) {
    const [isFetching, setIsFetching] = React.useState(false);

    React.useEffect(() => {
        setIsFetching(true)

        fetchFunction()
            .then(response => response.json())
            .then(data => {
                setIsFetching(false)
                dataHandler(data)
            })
            .catch(error => {
                setIsFetching(false)
                errorHandler(error)
            })
    }, []);

    return  isFetching;
}