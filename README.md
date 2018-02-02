# slicer

Fast 3D mesh slicer written in Go. Writes slices to grayscale PNG files.

### Install Go

First, install Go, set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```bash
brew install go
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

### Install Slicer

```bash
$ go get -u github.com/fogleman/slicer/cmd/slicer
```

### Example Usage

```bash
$ slicer --help

# slice model.stl with slices that are 0.1 units thick, rendering PNGs that
# cover 100x100 units in size with resolution of 10 pixels per unit
$ slicer -s 0.1 -w 100 -h 100 -x 10 model.stl 
```
