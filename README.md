# i3layout - Tools for i3wm

`i3layout` is a set of tools that extend i3wm, giving it support for dynamic
layouts and simple window movement in the style of xmonad.

See [this snippet](docs/i3config) of my i3 config for a complete usage example.

## i3layout

`i3layout` provides dynamic layouts.

By default, i3wm supports only very basic window layouts: splitting the screen
evenly in one direction, or showing one window at a time. If you want a 30/70
split, a master area, or any other kind of layout offered by other tiling window
managers such as xmonad, you have to do that manually.

`i3layout` solves this problem by running a server that listens for window
changes from i3 and lays out workspaces automatically as windows come and go. A
"oneshot" mode is also provided, allowing the user to trigger a manual re-layout
of the focused workspace.

### Installation

```console
$ go install github.com/adamwg/i3layout/cmd/i3layout
```

### Usage

To manually re-layout a workspace immediately, use the oneshot mode:

```console
$ i3layout oneshot [--layout=<name>] [workspace name]
```

By default, the `tall` layout will be applied to the focused workspace.

To have i3layout automatically manage your windows, use the server mode:

```console
$ i3layout serve
```

It is recommended that you run the server from your i3 config. It is *not* wise
to use the oneshot mode with the server running.

To change the layout of a workspace with the server running use the client:

```console
$ i3layout client change-layout [--workspace=<name>] <layout name>
```

By default, the focused workspace's layout is changed. The special names `prev`
and `next` can be used in place of a layout name to cycle through layouts.

### Supported Layouts

i3layout supports the following layouts:

#### `tall`

The tall layout has a single narrow "master" window on the left with the
remaining windows tiled vertically on the right using the remainder of the
screen.

```
----------------------------------------------------
| window 1      | window 2                         |
|               |                                  |
|               |                                  |
|               |----------------------------------|
|               | window 3                         |
|               |                                  |
|               |                                  |
|               |----------------------------------|
|               | window 4                         |
|               |                                  |
|               |                                  |
----------------------------------------------------
```

#### `columns`

This is i3's built-in SplitH layout.

```
----------------------------------------------------
| window 1       | window 2       | window 3       |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
|                |                |                |
----------------------------------------------------
```

#### `rows`

This is i3's bulit-in SplitV layout.

```
----------------------------------------------------
| window 1                                         |
|                                                  |
|                                                  |
|--------------------------------------------------|
| window 2                                         |
|                                                  |
|                                                  |
|--------------------------------------------------|
| window 3                                         |
|                                                  |
|                                                  |
----------------------------------------------------
```

#### `full`

This is i3's built-in Tabbed layout.

### How it works

`i3layout` (ab)uses i3's layout saving feature. To lay out a workspace, it does
the following:

1. Lists all the windows on the workspace (in depth-first order, as described
   below).
2. Creates a "template" of window containers based on the number of windows and
   the current layout mode. Each container has an i3 "mark" based on the window
   ID of a window in the list.
3. Creates a new, temporary workspace containing empty containers based on the
   template.
4. Swaps each window with its marked container.
5. Swaps the temporary workspace and the real workspace by renaming them.
6. Kills the temporary workspace.

Since this is a somewhat involved process, i3layout users may notice a bit of
flicker when layout is happening.

## i3window

`i3window` implements dimensionless window navigation and movement and provides
some helpful tools for scripting and interrogating i3.

By default, i3 supports navigating and swapping windows in four dimensions:
left, right, up, and down. Other window managers "flatten" the window tree,
allowing the user to navigate or swap the "next" or "previous" window using some
arbitrary ordering.

`i3window` solves this by flattening the tree on the fly and allowing
next/previous movement and swapping. It also provides some helpful utilities for
scripting i3 and observing its state.

### Installation

```console
$ go install github.com/adamwg/i3layout/cmd/i3window
```

### Usage

`i3window` never moves focus or windows off the current workspace or current
output. It enumerates windows in a depth-first manner. E.g., given the window
tree:
```
                  workspace
                  /       \
               win1      splitv
                         /    \
                      splith  win2
                      /    \
                    win3  win4
```
The order of windows will be:
```
win1
win3
win4
win2
```

#### Navigation

To change focus, use `i3window focus next` or `i3window focus prev`.

#### Swapping

To swap the current window with the next or previous window use `i3window swap
next` or `i3window swap prev`.

#### Utilities

To get a json-formatted list of windows on a workspace, in the order i3window
enumerates them, run `i3window list [--workspace=<workspace name>]`. By default
windows from the focused workspace will be listed.

To get a list of all windows, use `i3window list --all`.

To get a json representation of i3's window tree, use `i3window tree`.

## i3listen

`i3listen` is a simple utility that listens for i3 events and prints them in
json format. It is helpful when developing against the i3 API.

```console
$ go get github.com/adamwg/i3layout/cmd/i3listen
$ i3listen window
{"change":"focus","container":{"id":94737508960384,"name":"emacs@gitega","type":"con","border":"normal","current_border_width":2,"layout":"splith","percent":0.7,"rect":{"x":576,"y":53,"width":1344,"height":1027},"window_rect":{"x":2,"y":0,"width":1340,"height":1025},"deco_rect":{"x":576,"y":0,"width":1344,"height":26},"geometry":{"x":0,"y":0,"width":896,"height":864},"window":65174629,"window_properties":{"title":"emacs@gitega","instance":"emacs","class":"Emacs","window_role":"","transient_for":0},"urgent":false,"focused":true,"focus":[],"nodes":[],"floating_nodes":[]}}
{"change":"focus","container":{"id":94737508011008,"name":"dev | i3layout","type":"con","border":"normal","current_border_width":2,"layout":"splith","percent":0.3,"rect":{"x":0,"y":53,"width":576,"height":1027},"window_rect":{"x":2,"y":0,"width":572,"height":1025},"deco_rect":{"x":0,"y":0,"width":576,"height":26},"geometry":{"x":0,"y":0,"width":884,"height":580},"window":58720260,"window_properties":{"title":"dev | i3layout","instance":"st-256color","class":"st-256color","window_role":"","transient_for":0},"urgent":false,"focused":true,"focus":[],"nodes":[],"floating_nodes":[]}}
```

## TODO / Help Wanted

`i3layout` does what I need, but could do more. Any contributions would be
welcome, but especially the following:

* More layouts.
* File-based configuration of layouts.
* Reduced code duplication between layouts.
* Bug fixes.
* Fixing any incompatibilities with [x3](https://github.com/eonpatapon/x3) since
  they're great to use together.
