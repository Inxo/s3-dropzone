# File Drop App

File Drop App is a simple Golang application with a graphical interface using the Fyne library. The application provides a file drop zone, and after dragging a file into this zone, it displays the path to that file.

## Installation

Firstly, ensure that you have Golang and Fyne installed. Install Fyne using the following command:

```bash
go get fyne.io/fyne/v2
```

Then clone the repository and run the program:

```
git clone https://github.com/inxo/s3-dropzone.git
cd file-drop-app
go run main.go
```

## Usage

Run the application by executing go run main.go.
The "File Drop App" window will appear with a drop zone and the label "Drop File Here."
Drag a file into the drop zone.
Afterward, the path to the selected file will be displayed below the drop zone.

## Dependencies

Fyne - A UI toolkit and app API written in Go.

## License

This project is licensed under the [MIT License](LICENCE.md).