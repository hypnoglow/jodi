# jodi

A simple job dispatcher.

A library implementation of a 2-tier channel system pattern.  
Original idea by [Marcio Castilho](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/).

## Purpose

The example usage can be to offload any `http.Handler` from jobs 
that are waiting for an available worker and blocking the code.

See `example/main.go` for details.

## License

[MIT](https://github.com/hypnoglow/jodi/blob/master/LICENSE)