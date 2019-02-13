# Drone FPM Plugin
[![Build Status](https://cloud.drone.io/api/badges/lodge93/drone-fpm/status.svg)](https://cloud.drone.io/lodge93/drone-fpm)
[![GoDoc](https://godoc.org/github.com/lodge93/drone-fpm?status.svg)](https://godoc.org/github.com/lodge93/drone-fpm)

[Drone](https://drone.io/) plugin for [fpm](https://github.com/jordansissel/fpm).

This project generates an entrypoint script for drone based upon `fpm -h` output to combine repeatable, simple, 
configuration as code from drone with the ease and versatility of package management from fpm.

My personal use case is that I want to be able to define Debian package creation in code for my projects one time, and
then copy and paste that to all future projects that require Debian packaging 👌.

## Usage
This plugin supports all long hand fpm flags. All flag options can be used as plugin settings by substituting all `-`
with `_`. For example, the output type fpm option looks like this:
```bash
-t, --output-type OUTPUT_TYPE the type of package you want to create (deb, rpm, solaris, etc)
```
For this option, the shorthand `-t` is ignored, and the drone plugin accepts `output_type`.

For flags that do not have input (`--verbose` as an example), the drone plugin ignores the input given and simply passes
the flag to fpm, so any value would do (`verbose: true` as an example).

There is one special option for the drone plugin named `command_arguments` that is used for the positional arguments
required for the input type in use. For more clarity, you can see how this is used in the
[entrypoint template](assets/entrypoint.sh.templ).

For a complete list of options, see the [fpm documentation](https://fpm.readthedocs.io/en/latest/)

## Example
This example creates a Debian package with a systemd unit file for a binary named foo:

```yaml
- name: build debian package
  image: quay.io/lodge93/drone-fpm:latest
  settings:
    name: foo
    version: v0.0.1
    input_type: dir
    output_type: deb
    package: out/foo-v0.0.1.deb
    deb_systemd: foo.service
    command_arguments: out/foo=/usr/local/bin/
```

## Shell
This example shows how to run this plugin via `docker run`:

```bash
docker run --rm -v $(PWD):/workdir -w /workdir -e PLUGIN_DEB_SYSTEMD=/workdir/test/generate-entrypoint.service -e PLUGIN_NAME=generate-entrypoint -e PLUGIN_VERSION=snapshot-$(git log -n 1 --pretty=format:"%H") -e PLUGIN_INPUT_TYPE=dir -e PLUGIN_OUTPUT_TYPE=deb -e PLUGIN_PACKAGE=/workdir/out/generate-entrypoint-snapshot-$(git log -n 1 --pretty=format:"%H").deb -e PLUGIN_COMMAND_ARGUMENTS=/workdir/out/generate-entrypoint=/usr/local/bin/ quay.io/lodge93/drone-fpm:$(git log -n 1 --pretty=format:"%H")
```

## Tests
To run unit tests, run the following:
```bash
make test
```

This project also packages the entrypoint generation binary as a very basic end to end test. To test this project end
to end, run the following:
```bash
make end-to-end-test
```
