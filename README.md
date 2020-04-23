# ZipStream

A quick piece of experimental code to see whether it's possible to:

 - Stream a zip file from some Reader(ish) interface
 - Stream files within that zip
 - Perform some action on each item within each file within the zip
    - For csv files, an item is a line
    
## cmd/csvgen

```
Usage:
    csvgen <filename> <number of columns> <number of rows>
```
This tool will generate a CSV file with the specified filename that 
has the specified number of columns and rows.

It can be used to generate test CSVs for the `zs` tool

## cmd/zs
```
Usage:
    zs <filename>
```

The purpose of this tool is to demonstrate how we would stream zip files.  
It uses the zip on 
 