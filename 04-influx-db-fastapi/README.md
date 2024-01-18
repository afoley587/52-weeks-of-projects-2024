# Recording Time Series Data With APIs: InfluxDB + FastAPI
## Building an API that can read, query, and write data into FastAPI

## Introduction
Time series data has seen an explosion over the past few years, and for good
reason. It can be used everywhere. Look at Datadog, for example. Datadog
took time series data to the next level by building monitors, metrics, alerts,
etc. on top of that data. Tableau recently got into the time series data game
as well to provide better reporting capabilities to users. But, at its core, 
what is time series data?

## Time Series Data
Time series data is data that is recorded over consistent time intervals. It
can be used to see how certain data points change over time. For example, in our
project today, we will be monitoring heights of waves over some time period.
This wouldn't make too much sense in a typical relational database. What would
our tables look like? One table for locations, one for wave heights, and a relational
table to tie them together? How often do we insert rows into the wave heights table?
I think its clear to see we would have a large mess of sloppy data which would
make queries slow, complex, and hard to understand. This is where time series
data thrives! Now, time series data is structured and queried differently, so
something like postgres wouldn't be a natural fit. We can't just shoehorn time 
series data into a relational database engine. Enter [InfluxDB](https://www.influxdata.com/) 
which has been specifically built to handle time series data. Now, there are other time
series database engines (like TimeScale and CrateDB), but I like Influx...
so that's what we will use today!

## Project Overview And Setup
So, what are we building today? Well, we want to build an API that can read, write, 
and query data from InfluxDB. The database should allow users to:

1. Record a new wave height in some location
2. List all of the known waves and their locations
3. Filter the records based on location
4. Filter the records based on minimum wave height
5. Filter the records based on a combination of 3 and 4

We then want the API to serialize this data in a nice and predefined format and 
return it to some caller. We could imagine that the caller would be some charting
frontend that tells surfers where big swells have been over the course of the past
day or week or month.

For our API, we will be using Python and FastAPI. We will add two routers: one 
to read data and one to write data. We will then implement an InfluxDB client
which somewhat restricts the overall power of the 
[influx-python-client](https://github.com/influxdata/influxdb-client-python/tree/master)
library.

The directory structure can be seen below:

```shell
.
├── __init__.py
├── client
│   ├── __init__.py
│   └── influx.py
├── config.py
├── main.py
├── routes
│   ├── __init__.py
│   ├── read.py
│   └── write.py
└── schemas.py
```

We see that we have a `client` directory, which is where we will
write our restricted client library. We also have a `routes`
directory whcih will house our read/write API routers. And then
we have a few files: `config.py`, `main.py`, and `schemas.py`.
`config.py` is where we will place our InfluxDB connection  settings.
`main.py` is where we will actually start/run our API and glue
everything together. And finally, `schemas.py` is where we will put
our pydantic models that our API and client will use.

There are a few external libraries to install, such as `influx-python-client`
and `fastapi`. I have provided a [requirements.txt]() for ease of use!

## Building The API
Now that we know what we are building, let's begin building it. I have
broken down the process into 5 steps:

1. Writing Our Pydantic Models
2. Writing Our Client
3. Writing Our Write Router
4. Writing Our Read Router
5. Configuration And Tying It Together

So, let's begin!

### Step 1: Writing Our Pydantic Models
### Step 2: Writing Our Client
### Step 3: Writing Our Write Router
### Step 4: Writing Our Read Router
### Step 5: Configuration And Tying It Together

## Running The Stack
### Step 1: Bootstrapping InfluxDB
### Step 2: Using Our API