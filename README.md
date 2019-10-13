# jokesontap
Have you ever wanted to query Chuck Norris-like jokes as fast as possible?  Of course you have.  Get them while
they're hot.

This application server retrieves fresh jokes from the [Internet Chuck Norris Database](http://www.icndb.com/), but
adds a flare of personality by switching out the name with a random one from [uinames.com](https://uinames.com/).
Why should Chuck get all the credit?

## Building
`build.sh` will test and build the server binary.
```bash
cd jokesontap
./build.sh
ls ./bin
```

## Usage

### Starting the Server
After [building](#building) the binary can be run with defaults.
```bash
chmod +x ./bin/jokesontap
./bin/jokesontap 
```

Or set server options.
```bash
./bin/jokesontap --help
./bin/jokesontap --port 8080 --log-level error
```

### Querying
The server has a single root endpoint which will return a new Chuck Norris-like joke with a random name.

Assuming the server is running on default port 5000, query the server and get a joke.
```bash
$ curl http://localhost:5000
John Smith's OSI network model has only one layer - Physical.
```

## Known Limitations
As of writing [uinames.com](https://uinames.com/), which is used to generate the random names, has a rate limit of 3,500
names per minute.  This is partially mitigated by eagerly querying and storing names in memory, but if pushed the server
may not be able to serve a new joke for lack of a random name.

## TODO
- [ ] Implement caching so that when given the Cache-Control header the server will reuse a previous name.
- [ ] Implement Prometheus metrics.
