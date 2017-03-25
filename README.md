# bindist

[![Build Status](https://travis-ci.org/ross-spencer/bindist.svg?branch=master)](https://travis-ci.org/ross-spencer/bindist)
[![GoDoc](https://godoc.org/github.com/ross-spencer/bindist?status.svg)](https://godoc.org/github.com/ross-spencer/bindist)
[![Go Report Card](https://goreportcard.com/badge/github.com/ross-spencer/bindist)](https://goreportcard.com/report/github.com/ross-spencer/bindist)

Calculate distance between two hexadecimal values in a given file. 

Example use case is file format signature development, e.g. https://groups.google.com/forum/#!topic/droid-list/QRQ9LScT8iw 

      Usage:  bindist [-magic1 ...] [-magic2 ...] [-file ...]
                      [Optional -size] [Optional -fname]
                      [Optional -version]

      Output: [CSV] 'offset','size','filename'
      Output: [STRING] 'bindist version ...'                     

        -file string
          	File to find the distance between.
        -fname
          	[Optional] Return filename alongside offset and size.
        -magic1 string
          	First magic number in a file to begin from, and offset, e.g. magic,offset.
        -magic2 string
          	Second magic number in a file to search for, no offset, e.g. magic.
        -size
          	[Optional] Return size of file alongsize offset in bytes.
        -version
            [Optional] Return version of bindist.

### Todo

- Potentially use a graphics library to scatter plot for us.
- Attempt to enable Linux globbing / stdio

## License

**[GPL Version 3](http://choosealicense.com/licenses/gpl-3.0/)**: https://github.com/ross-spencer/bindist/blob/master/LICENSE
