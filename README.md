# washeet
Washeet is a spreadsheet web ui package written in golang. The project aims to provide a fast, light-weight, well polished spreadsheet ui targeting web applications that need to present streaming tabular data and allow user interactions. As of now it is very much in its early stages and only has some basic features like -

* setting cell contents via api, but cells not yet user editable.
* control cell font size, attributes like bold, italics and underline via api
* control cell foreground/background via api
* cell content alignment control via api.
* all above features can be changed any time while sheet is running, taking advantage of 60 fps canvas.
* individual column-width and row-height control via api.
* cell cursor movements with mouse as well as with keyboard.
* cell range selection with mouse as well as with keyboard.
* copying cell/cell-range contents to clipboard.

## API documentation
Things are still in flux, but see [godoc](https://godoc.org/github.com/dennisfrancis/washeet) for the current form.

## Is there a demo/example ?
Yes, there is one in the `demo` sub-directory. The demo shows how to feed your data to the spreadsheet. The data in the demo is just random numbers but it could potentially come from outside, say from a websocket connection(s) from the server. See below for how to setup the demo. You must have installed Go version >= 1.11 and "GNU make" to do this. Technically you don't need "make" but it would certainly make things easier if you do. Obviously you would need a very recent browser with good web-assembly support like chrome/firefox. 

```
# Get the package.
$ go get -u github.com/dennisfrancis/washeet

# Go to package's src root.
$ cd $GOPATH/src/github.com/dennisfrancis/washeet

# Build the demo
$ make demo
```

Now run an https server on demo sub-directory. Any https server would do. If you use just http, the clipboard api wont work. You could try [caddy](https://caddyserver.com/download) as it is easy to setup without writing a single line of code.
```
$ cd demo
$ caddy  # assuming you chose to use caddy https server.
```
Open Chrome browser and point to "https://127.0.0.1:1234". You should see a spreadsheet ui with some contents on it and now you can play around with it!

## Contributing
You are always welcome to submit a bug-report(issues) here or provide a pull request. Happy hacking !
