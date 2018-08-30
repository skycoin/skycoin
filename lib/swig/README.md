#SWIG files

SWIG interface files for generating skycoin libraries in different languages.

##Table of Contents

<!-- MarkdownTOC levels="1,2" autolink="true" bracket="round" -->
- [Requirements] (#requirements)
- [Usage] (#usage)
	- [Tips] (#tips)


## Requirements
  Requires Swig installed. It has been tested with 3.0.10, so that would be the preferred version.
  These swig interface files has been tested with Python2.7, and in other languages may not work as expected. It is very important to read SWIG documentation for the specific target language being used, http://www.swig.org/Doc3.0/SWIGDocumentation.html
## Usage	
  First step would be to apply SWIG to the main interface file (skycoin.i). For example if it is going to be generated for Python it would be like this:
  swig -python skycoin.i
Or in the case of C#:
  swig -csharp skycoin.i
  You will also need to add to this command the include paths where to search for required header files. This path would be skycoin include path. For example if skycoin include path is skycoin/include, then the swig command would be:
  swig -csharp -Iskycoin/include skycoin.i
  However, doing this will raise an error because SWIG doesn't fully understand a specific line in libksycoin.h. The solution is making a copy of libskycoin and removing the conflicting line. Like this:
  grep -v _Complex skycoin/include/libskycoin.h > swig/include/libskycoin.h
  swig -csharp -I swig/include -Iskycoin/include skycoin.i
  The above command lines will make a copy of libskycoin.h removing lines containing string "_Complex" and execute swig specifying first as include path the path containing the modified copy of libskycoin.
## Tips
  It is also important to notice that file structs.i contains inclusions of header files containing type definitions of libskycoin. This file is very similar to skytypes.gen.h but replacing character # with %.
  For example skytypes.gen.h being something like this:
  
  #include "file1.h"
  #include "file2.h"
  
  Then, structs.i would be:
  %include "file1.h"
  %include "file2.h"
  
  And this is something that swig can understand.
  So, a good idea would be to update structs.i from skytypes.gen.h, just in case there have been modifications to skycoin source code.
  Like this:
  cp skycoin/include/skytypes.gen.h structs.i
  sed -i 's/#/%/g' structs.i
  grep -v _Complex skycoin/include/libskycoin.h > swig/include/libskycoin.h
  swig -csharp -I swig/include -Iskycoin/include skycoin.i
  
  Use https://github.com/simelo/pyskycoin/ as a reference. Check Makefile rule build-swig.
  This project is used to create a Python extension to access skycoin API from Python.  This project contains skycoin repository as git submodule, which is a good idea to make libraries for other languages.
  
  
  
  
    
