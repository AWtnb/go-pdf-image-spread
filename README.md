# README

Convert PDF to image file and concatenate side by side as in book spread.

## Build

```
go build -o converter.exe main.go
```

## Usage

+ Converts PDF to image in the same folder as `converter.exe` .
+ Create a folder with the same name as the target PDF and store the image file of each page in it.
    + Create a subfolder named `conc`, and output the combined final file in it.
+ If the folder already exists, an error will occur and the process will be aborted.

### With mouse

Place the `.bat` file you want to use in the same folder as `converter.exe` and double-click it.

+ `horizontal.bat`
+ `horizontal-singletop.bat` (Horizontal PDF. single first page.)
+ `vertical.bat`
+ `vertical-singletop.bat` (Vertical PDF. single first page.)

### From command line

Command line options:

+ `--singleTop` : If the PDF has an odd number of pages, leave the first file on a single-page.
    + By default, leave the last page single-page.
+ `--vertical` : When concatenating vertical PDFs, such as Japanese, the pages are aligned from right to left.

