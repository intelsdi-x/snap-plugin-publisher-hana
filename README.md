# Snap publisher plugin - SAP Hana

This plugin supports pushing metrics into the SAP Hana server.

It's used in the [Snap framework](http://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

* [Snap](https://github.com/intelsdi-x/snap)
* [golang 1.6+](https://golang.org/dl/) (needed only for building)
* [SAP Hana server](http://hana.sap.com/abouthana.html)

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64
* Darwin/amd64

### Support Matrix

- Hana Plugin: v4 -> Snap version 0.9.0-beta

### Installation

#### Download hana plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-publisher-hana/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-publisher-hana

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-hana.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* As part of plugin configuration in Task manifest provide:
  * `host` default `root`
  * `port` default `30017`
  * `username` default `root`
  * `password` default `root`
  * `database` default `SNAP_TEST`
  * `tablename` default `info`

## Documentation
There are a number of other resources you can review to learn to use this plugin:

* [Snap hana integration test](https://github.com/intelsdi-x/snap-plugin-publisher-hana/blob/master/hana/hana_integration_test.go)
* [Snap hana unit test](https://github.com/intelsdi-x/snap-plugin-publisher-hana/blob/master/hana/hana_test.go)

### Roadmap

There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-hana/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-hana/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions! 

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Jon Machen](https://github.com/jkmachen)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
