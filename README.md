# bindist

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

###TODO

- Potentially use a graphics library to scatter plot for us.
- Attempt to enable Linux globbing

##License

Copyright (c) 2016 Ross Spencer

This software is provided 'as-is', without any express or implied warranty. 
In no event will the authors be held liable for any damages arising from the 
use of this software.

Permission is granted to anyone to use this software for any purpose, including 
commercial applications, and to alter it and redistribute it freely, subject to 
the following restrictions:

- The origin of this software must not be misrepresented; you must not claim that 
you wrote the original software. If you use this software in a product, an 
acknowledgment in the product documentation would be appreciated but is not 
required.

- Altered source versions must be plainly marked as such, and must not be 
misrepresented as being the original software.

- This notice may not be removed or altered from any source distribution.
