# UwU Note

| OS | Build-Status |
| - |:- |
| linux | [![CircleCI - Linux](https://circleci.com/gh/UwUNote/uwunote/tree/master.svg?style=svg)](https://circleci.com/gh/UwUNote/uwunote/tree/master) |
| darwin | [![TravisCI - Mac OS](https://travis-ci.org/UwUNote/uwunote.svg?branch=master)](https://travis-ci.org/UwUNote/uwunote) |
| windows | [![AppVeyor - Windows](https://ci.appveyor.com/api/projects/status/7e3yaftwricytyj5/branch/master?svg=true)](https://ci.appveyor.com/project/Bios-Marcel/uwunote/branch/master) |

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

## Installation using pre-built binaries

The Linux binary can be found [here](https://circleci.com/gh/UwUNote/uwunote/). Just click the latest **green**(passing) build and download it under artifacts.

A Windows executable without dependencies and an installer with dependencies can be found [here](https://ci.appveyor.com/project/Bios-Marcel/uwunote/branch/master/artifacts).
If you want to manually install the dependencies, check [this tutorial](https://gianlucacosta.info/go-gui-apps-on-windows) out.

MacOS binaries are currently not available.

## Building manually

You need the development packages of following libraries:

* GTK+3 (3.22 and later)
* GDK 3
* GLib 2
* Cairo

To build the app you simply run these commands:

```bash
cd $GOPATH
go get github.com/UwUNote/uwunote
cd src/github.com/UwUNote/uwunote
go get ./...
go build
```

In case you want to put the binary into your go path, run

```bash
go install
```

as the last command instead, that will put a ready to use binary called `uwunote` into `$GOPATH/bin`.

### Linux

You can most likely install all dependencies via your systems packagemanager.

### Windows

Follow the tutorial mentioned in the comment and apply the fix manually.

[https://github.com/UwUNote/uwunote/issues/1#issuecomment-414850043](https://github.com/UwUNote/uwunote/issues/1#issuecomment-414850043)

## Configuration

Configuration is currently done by editing a JSON file, check the wiki for more information:
[https://github.com/UwUNote/uwunote/wiki/Configuration](https://github.com/UwUNote/uwunote/wiki/Configuration)

In case the documentation is outdated, you can always check the applications settings dialog. The settings dialog contains tooltips for every setting on each label.

## Contribute

Feel free to create issues or create pull requests. Issues may also request features or just be questions.