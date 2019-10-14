# jokesontap
Have you ever wanted to query Chuck Norris-like jokes as fast as possible?  Of course you have.  Get them while
they're hot.

This application server retrieves fresh, specifically nerdy, jokes from the
[Internet Chuck Norris Database](http://www.icndb.com/), but adds a flare of personality by switching out the
name with a random one from [uinames.com](https://uinames.com/). Why should Chuck get all the credit?

## Building
`build.sh` will test and build the server binary.
```bash
cd jokesontap
./build.sh
ls ./bin
```

## Features
- fast, concurrent web server
- ahead-of-time random name cache partially mitigates backpressure from names API and avoids API rate limiting
- detailed logging
- customized application settings through command line parameters
- automatic build and testing through build script
- well tested code, of course

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
Bruce Banner's OSI network model has only one layer - Physical.
```

## Known Limitations
As of writing [uinames.com](https://uinames.com/), which is used to generate the random names, has a rate limit after
a certain number of requests.  This is partially mitigated by eagerly querying and storing names in memory, but
if pushed the server may not be able to serve a new joke for lack of a random name.

## Benchmarks
Benchmarking the server for 30 seconds, after the names cache (10,000 entries) was allowed to fill.  Disclaimer:
benchmark results were run with the server and benchmarking client on the same laptop.

The results are as follows:
```
$ timeout 60 siege -b -c 100 http://localhost:5000/
~~~
Transactions:                  13000 hits
Availability:                  97.01 %
Elapsed time:                  59.99 secs
Data transferred:               0.96 MB
Response time:                  0.43 secs
Transaction rate:             216.70 trans/sec
Throughput:                     0.02 MB/sec
Concurrency:                   92.77
Successful transactions:       13000
Failed transactions:             400
Longest transaction:            5.30
Shortest transaction:           0.15
```

Ultimately this service's current bottleneck is the name server which we depend on to generate random names.  The
name server throttles our traffic and thus we cannot generate names fast enough.  This could be remedied by some of the
enhancements in [TODO](#todo).

## TODO
- [ ] Implement caching so that when given the Cache-Control header the server will reuse a previous name always,
or when the name server is unavailable.
- [ ] Implement Prometheus metrics.
- [ ] As of writing, the [Internet Chuck Norris Database](http://www.icndb.com/) only contains 574 jokes total.  This
is such a small database it would be more efficient to simply download the entire data set on a daily basis and serve
the jokes from memory without calling out to the external service.  This functionality needs to be confirmed with the
specific application requirements before implementation.
