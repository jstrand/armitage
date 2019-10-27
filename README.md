Armitage
===============================================================================

Experimental file copy tool in Go.

* Make backups of new files from given directory.
* File copies are put in a directory called <files>.
* A folder called <checksums> contains checksums of already copied files.

Changed or previously copied files are ignored (useful for copying photos/video).
Integrity of files can be checked with ```shasum --check ../checksums/*``` from the files folder

Usage
-------------------------------------------------------------------------------
armitage <path-to-copy-from>
