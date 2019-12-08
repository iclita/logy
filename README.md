# Logy

The best parser for filtering and handling files of any size with ease. Inspecting or filtering large files can be a real pain. Logy let's you open any file in a paginated manner so that there is no overhead of opening the whole file when all you want is searching for a small chunk of text. Besides file paths, it also supports folder paths.

## Features

- Written in pure Go (Golang)
- Requires Go >= 1.11 to build. Just visit [this link](https://golang.org/doc/install) on how to install Go on your machine
- Simple installation process
- Very easy to use and highly intuitive
- Works on Mac, Linux and Windows

## Installation

There is no need to install anything. If you want to skip the setup process and just use the tool, follow [this link](https://github.com/iulianclita/logy/releases) and use the build that is specific to your OS (Mac, Linux or Windows) and architecture (32 or 64 bit)

Otherwise clone the repository. This way you can contribute to the project :)

## Usage

After the file is parsed you can navigate to any page by specifying the desired page number

### Basic usage
```bash
$ logy path/to/file.log # By default it outputs 50 lines per page
```

```bash
$ logy path/to/folder --ext=log,txt # The extensions (-ext) flag must be specified for folder paths to mention what file types should be scanned. In this example the parser will search the folder recursively for all files ending with .log or .txt extension
```

### Specify how many lines per page
```bash
$ logy path/to/file.log --lines=25 # Now it will output 25 lines per page
```

```bash
$ logy path/to/folder --ext=log,txt --lines=25 # Now it will output 25 lines per page
```

### Format json output
```bash
$ logy path/to/file.log --text=json # Every json structure that is found will be nicely formatted 
```

```bash
$ logy path/to/folder --ext=log,txt --text=json # Every json structure that is found will be nicely formatted 
```

### Search for text
```bash
$ logy path/to/file.log --filter=Exception # Every text that is found will be nicely colored to be easily observed 
``` 

```bash
$ logy path/to/folder --ext=log,txt --filter=Exception # Every text that is found will be nicely colored to be easily observed 
``` 

### Navigate to any page
```bash
$ logy path/to/file.log --page=10 # The parser will directly navigate to the specified page number 
```

```bash
$ logy path/to/folder --ext=log,txt --page=10 # The parser will directly navigate to the specified page number 
```

### Enable regex support
```bash
$ logy path/to/file.log --filter=[0-9]{2}:[0-9]{2}:[0-9]{2} --with-regex # The parser will search for any text that matches whatever was specified in the filter option flag
```

```bash
$ logy path/to/folder --ext=log,txt --filter=[0-9]{2}:[0-9]{2}:[0-9]{2} --with-regex # The parser will search for any text that matches whatever was specified in the filter option flag
```   

### Disable colored output of any kind
```bash
$ logy path/to/file.log --no-color # The parser will display all text with the same color (black/white). Probably you will never want this behavior but it's here just in case :)
``` 

```bash
$ logy path/to/folder --ext=log,txt --no-color # The parser will display all text with the same color (black/white). Probably you will never want this behavior but it's here just in case :)
``` 

Of course all the flag options can be combined in any manner to obtain the desired results

## Note
Because regex implementation in Go is not highly performant, use the `--with-regex` flag when it is absolutely necessary, especially with large files.

## Contributing

#### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/iulianclita/logy/issues) to report any bugs or feature requests.

## License

The Logy CLI tool is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT).
