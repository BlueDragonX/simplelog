Go Simple Logger
================
A simple logger that wraps log and log/syslog into something consistent and
easy to use.

Example
-------
Declare the logger:

    logger, _ := simplelog.NewLogger(simplelog.CONSOLE | simplelog.SYSLOG, "example")

Send logs to it:

    logger.Notice("starting the app")

Close the logger:

    logger.Close()

License
-------
Copyright (c) 2014 Ryan Bourgeois. Licensed under BSD-Modified. See the LICENSE
file for a copy of the license.
