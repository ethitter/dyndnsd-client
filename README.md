dyndnsd-client
==============

Client for the [`dyndnsd`](https://github.com/cmur2/dyndnsd) daemon. Set up `dyndnsd` first, otherwise this is useless.

## Configuration

Copy `config-sample.json` to `config.json` and update the values.

## Use

Set up a cron entry to run the program periodically, or manually call it when your IP address changes. The program accepts no arguments, as `config.json` handles everything.
