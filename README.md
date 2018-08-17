# UwU Note

First of all, I really hate making up names, so `UwU Note` is all you get.

## What's it?

It's a simple note-taking app. Notes stick to your desktop, you can create new ones and delete existing ones.
Every note has it's own window, which saves it's coordinates and size.

## The Future

For the future i have planned some stuff:

* Cloud synchronization (Different Services)
* An android app
* Markdown support

## Installation

As this is a work in progress, there are currently no installation instructions, just build the project yourself instead.
Theretically you should be able to build it on every platform supporting GTk3+.

Official binaries are gonna follow at some point, maybe even a snap or flatpak.

## Building

To build the app you simply run these commands:

```bash
cd $GOPATH
go get github.com/Bios-Marcel/uwuNote
cd src/github.com/Bios-Marcel/uwuNote
go get ./...
go build
```

You need the development packages of following libraries:
* GTK+3 (3.12 and later)
* GDK 3 (3.12 and later)
* GLib 2 (2.36 and later)
* Cairo (1.10 and later)

You can most likely install those via your systems packagemanager.

## Configuration

Configuration is currently done by editing a JSON file, check the wiki for more information:
https://github.com/Bios-Marcel/uwuNote/wiki/Configuration

## Contribute

Feel free to create issues or create pull requests.