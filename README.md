# UwU Note

[![Waffle.io - Columns and their card count](https://badge.waffle.io/UwUNote/uwunote.svg?columns=all)](https://waffle.io/UwUNote/uwunote)

| linux | [![CircleCI - Linux](https://circleci.com/gh/UwUNote/uwunote/tree/master.svg?style=svg)](https://circleci.com/gh/UwUNote/uwunote/tree/master) |
| - | - |
| **darwin** | [![TravisCI - Mac OS](https://travis-ci.org/UwUNote/uwunote.svg?branch=master)](https://travis-ci.org/UwUNote/uwunote) |
| **windows** | **TODO** |

First of all, I really hate making up names, so `UwU Note` is all you get.

## What's it

It's a simple note-taking app. Notes stick to your desktop, you can create new ones and delete existing ones.
Every note has it's own window, which saves it's coordinates and size.

![Demo Image](https://i.imgur.com/tM3fhoK.jpg)

## The Future

For the future i have planned some stuff:

* Cloud synchronization via different services
* [An android app](https://github.com/UwUNote/uwunote-android)
* Markdown support

## Installation

As this is a work in progress, there are currently no installation instructions, just build the project yourself instead.
Theretically you should be able to build it on every platform supporting GTk3+.

Official binaries are gonna follow at some point, maybe even a snap or flatpak.

As of now, you can do the followig:

```bash
cd $GOPATH
go get github.com/UwUNote/uwunote
cd src/github.com/UwUNote/uwunote
go get ./...
go install
```

which will put a ready to use binary called `uwunote` into `$GOPATH/bin`.

## Building

To build the app you simply run these commands:

```bash
cd $GOPATH
go get github.com/UwUNote/uwunote
cd src/github.com/UwUNote/uwunote
go get ./...
go build
```

You need the development packages of following libraries:

* GTK+3 (3.12 and later)
* GDK 3 (3.12 and later)
* GLib 2 (2.36 and later)
* Cairo (1.10 and later)

### Linux

In addition to the default dependencies, linux needs additional dependencies:

* libappindicator

You can most likely install all dependencies via your systems packagemanager.

### Mac OS

TODO

### Windows

[https://github.com/UwUNote/uwunote/issues/1#issuecomment-414850043](https://github.com/UwUNote/uwunote/issues/1#issuecomment-414850043)

## Configuration

Configuration is currently done by editing a JSON file, check the wiki for more information:
[https://github.com/UwUNote/uwunote/wiki/Configuration](https://github.com/UwUNote/uwunote/wiki/Configuration)

## Contribute

Feel free to create issues or create pull requests. Issues may also request features or just be questions.