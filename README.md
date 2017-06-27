# slicer

Fast 3D mesh slicer written in Go. Writes slices to an SVG file.

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
$ slicer -s 0.1 model.stl  # slice model.stl with slices that are 0.1 thick
$ slicer -n 100 model.stl  # slice model.stl into 100 slices
$ slicer -q -s 0.1 -d . folder/*.stl  # quietly slice multiple stl files and put results in cwd
```

### Example Output

![Example](http://i.imgur.com/BNKZ8HY.png)
