`dyndnsd-client`
================

Client for the [`dyndnsd`](https://github.com/cmur2/dyndnsd) daemon. Set up `dyndnsd` first, otherwise this is useless.

## Use

Set up a cron entry to run the program periodically, or manually call it when your IP address changes. The program accepts a single argument, `-config`, pointing to a `config.json` file that handles everything. A sample `config.json` is provided in the form of `config-sample.json`.

## Configuration Options

* `username`: daemon username
* `password`: daemon password
* `protocol`: `http` or `https`, preferably the latter
* `host`: IP or hostname of `dnydnsd` instance
* `port`: port `dyndnsd` is listening on
* `path`: path to `dnydnsd` daemon, typically `nic/update`
* `dns_hostname`: the `nsd` zone to update
* `ipv4_endpoint`: GET-accessible endpoint that returns an IPv4 IP address as a string; when in doubt, use `http://whatismyip.akamai.com/`
* `ipv6_endpoint`: GET-accessible endpoint that returns an IPv6 IP address as a string; when in doubt, use `http://ipv6.whatismyip.akamai.com/`
* `ipv6_use_prefix`: if running multiple instances on the same network, use only the IPv6 prefix, omitting the instance portion of the address; this greatly reduces unnecessary `nsd` updates

## Proxying via `nginx`

`dnydnsd` doesn't support HTTPS, so using `nginx` as an HTTPS proxy is recommended. Sending a username and password over an unsecured connection should be avoided if possible.

More coming soon...
