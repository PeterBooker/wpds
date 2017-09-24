[![Release](https://img.shields.io/badge/version-v0.4.0-blue.svg)](https://github.com/PeterBooker/WordPress-Directory-Slurper/releases/tag/v0.4.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/PeterBooker/WordPress-Directory-Slurper)](https://goreportcard.com/report/github.com/PeterBooker/WordPress-Directory-Slurper)
[![License](https://img.shields.io/badge/license-GPL--2.0%2B-red.svg)](https://github.com/PeterBooker/WordPress-Directory-Slurper/blob/master/LICENSE)

# WPDS (WordPress Directory Slurper)
WPDS is a cross-platform CLI tool built with [Go](https://golang.org/). Slurps down the latest version of every Plugin and/or Theme in the WordPress Directory. Comes with inbuilt searching and formatted search summaries. Based on the Plugin and Theme Directory Slurpers by [markjaquith](https://github.com/markjaquith/WordPress-Plugin-Directory-Slurper), [clorith](https://github.com/Clorith/WordPress-Plugin-Directory-Slurper), [ipstenu](https://github.com/Ipstenu/WordPress-Theme-Directory-Slurper) and [chriscct7](https://github.com/chriscct7/WordPress-Plugin-Directory-Slurper/).

Note: WPDS is still in early development and therefore may contain bugs or miss features.

## Dependencies

None. WPDS is a self-contained executable.

## Install

Download the relevant file for your operating system from the [releases](https://github.com/PeterBooker/WordPress-Directory-Slurper/releases) page, then either run it from the directory you want it to work in or put it into your PATH and it will use the current working directory.

## Examples

### Download Plugin Directory

```
wpds download plugins
```

This will download the entire plugin directory.

### Update Plugin Directory

```
wpds update plugins
```

This will update the existing files using the latest revision found in `/plugins/.last-revision`.

## Features

- [x] Download the Plugin Directory
- [x] Update the Plugin Directory files
- [x] Download the Theme Directory
- [x] Update the Theme Directory files
- [x] Desktop Notifications
- [ ] In-built Searching
- [ ] Search Summary Generation

## FAQ

### Why did you remake the previous tools in Go?

Building the CLI tool in Go removes any requirements and provides full cross platform support, making it easier for everyone to use.

It also allowed me to build the search functionality into the tool, removing further requirements.

### Why download the zip files? Why not use SVN?

An SVN checkout of the entire repository is a BEAST of a thing. You don't want it, 
trust me. Updates and cleanups can take **hours** or even **days** to complete.

### How long will it take?

Your first update will take a while (at last testing around 2-3 hours using the default settings).

But subsequent updates are smarter. The script tracks the SVN revision number of your latest update and then asks the plugins SVN repository for a list of plugins that have changed since. Only those changed plugins are updated after the initial sync.

### How much disk space do I need?

As of September 2017:

69607 Plugins- extracting 884,056 folders, 2,321,090 files, taking up 34.6 GB of space.

### Something went wrong, how do I do a partial update?

The last successful update revision number is stored in `plugins/.last-revision`.
You can just overwrite that and the next `wpds update plugins` will start after that revision.

### What is this thing actually doing to my computer?

Once downloads have started, the CLI tool will display its progress including how many of the themes/plugins it has downloaded out of the total and an estimated completion time.

## License

GNU General Public License 2.0, see [LICENSE](https://github.com/PeterBooker/WordPress-Directory-Slurper/blob/master/LICENSE).