# GoCF
A terminal-based CHelper-like wannabe project for Go(lang) enthusiasts

Motivation
----------

I wanted to try and learn a bit the [Go programming language](https://golang.org/). And based on my prior experience, a nice 
indicator of how comfortable you can become in a programming language, is using it to solve problems of algorithmic nature. 
A few online judges provide support for Go, my favourites being [Codeforces](http://codeforces.com/) and 
[Timus](http://acm.timus.ru/). When using Java for solving problems I got used to 
[CHelper](http://codeforces.com/blog/entry/3273) and so I decided to write a tool providing similar (although more limited) 
functionality, but terminal-based and focused on Go. So *here* it is!

Naming
------
**GoCF** stands for **Go****C**ode**F**orces. Lame....I know.

Installation
------------

Make sure you have Go installed and `$GOPATH` set ([see here](https://golang.org/doc/install)). Then run
```
go get github.com/ale64bit/gocf
go install github.com/ale64bit/gocf
```
This should build the binary in `$GOPATH/bin/gocf`. You can later add it to your `$PATH` for convenience. For the rest 
of this document, we will skip the `$GOPATH/bin/` prefix.

Usage
-----

If you run `gocf` without arguments, it will show a usage message with a short description of the supported actions. Also, 
the first time you run it, you will be asked for 3 parameters:

- **work file**: this is the file where the source code of your solution lives. If you are working in a problem, this is the file 
where you write the code, and the file you submit to the online judge as well. If you save your session to work in a 
different problem, your work file will be saved as well, and restored at a later point.
- **session directory**: this directory contains the information for the current problem you are solving. This includes time 
limit, memory limit, input file, output file, checker, sample tests and their answers. 
- **archive directory**: when you want to work in a different problem, you can archive the current session, create a new 
one and then restore any other session at a later point to continue your work.

Each problem you ever worked on (i.e. session) is identified by two parameters: contest id and task id. When you archive 
a session, it will be stored according to this parameters in the archive directory. For example, if you are working in a 
session with contest `"swerc/2010/practice"` and task `"test"`, then your session will be stored at `$ARCHIVE_DIR/swerc/2010/practice/test`.

Example
-------
Let's solve a very simple problem, for illustrative purposes: http://codeforces.com/problemset/problem/71/A

The first thing we will do, is create a new session in order to solve this problem:
```
~$ gocf create
```
**NOTE**: if it's the first time you create a session, you will be asked for the 3 parameters mentioned above.

Before creating the new session, you will be asked if you want to save the current session (type `N` to ignore it). After 
that, you will be asked for the parameters of the session:
```
Configuration file found. Loading...
Do you want to archive current session? (yN): N
Enter contest name [default=practice]: codeforces/71
Enter task name [default=task]: A
Enter input file name [default=*]: 
Enter output file name [default=*]: 
Enter time limit [default=1000]: 2000
Enter memory limit [default=67108864]: 
Enter task checker [default=*]: 
done
~$ 
```
As you can see, we have accepted the default values for input file (stdin), output file (stdout), memory limit (64MB) and 
checker (the default checker is a port of [lcmp](https://github.com/MikeMirzayanov/testlib/blob/master/checkers/lcmp.cpp) 
but you can plug an arbitrary checker here). The default input file, output file and checker parameters are shown as 
`*`. After this, you can check the current session details with:
```
~$ gocf ls
Configuration file found. Loading...
Session description:
  Contest:    codeforces/71
  Task:       A
  Input:      *
  Output:     *
  Time limit: 2000 [ms]
  Mem limit:  64 [MiB]
  Checker:    *

TESTS
-----------------------------------------------------------------------
~$ 
```
As you noticed, there are no tests. Let's add the example test from the problem:
```
~$ gocf add
Configuration file found. Loading...
Session description:
  Contest:    codeforces/71
  Task:       A
  Input:      *
  Output:     *
  Time limit: 2000 [ms]
  Mem limit:  64 [MiB]
  Checker:    *


Enter input:
4
word
localization
internationalization
pneumonoultramicroscopicsilicovolcanoconiosis

Enter answer [empty if unknown]:
word
l10n
i18n
p43s
Added test # 1
~$ 
```
To end a test, type the EOF signal in your terminal (for me in Ubuntu 14.04, this is Ctrl-D). After this, you can check 
that the test was added by listing the session again (`gocf ls`). Note that the tests are indexed starting from 1. If you 
want to remove a test, just run `gocf rm <ID>`. If there are any tests with larger ID, they will shift down (that is, if 
there are 3 tests and you run `gocf rm 2`, the test that was #3 before, is not #2.

As you may have noticed, all this process of creating the session and entering the tests can become cumbersome. For 
Codeforces, you can reduce all those operations in a single one:
```
gocf import http://codeforces.com/problemset/problem/71/A
```
This will take care of creating a session with the proper details and also importing the sample tests from the problem 
statement. Right now, the import command only works for Codeforces, but let me know if there are other judges you would 
like to support, or even better, just send a PR with the change! :)

Now that we have the session created and the tests imported, we can start solving the problem. If you open the work file, 
you will find out that it contains already boilerplate for IO (according to the session input/output file specs). If you 
want to test the solution right now:
```
~$ gocf test
Configuration file found. Loading...
Removing test directory...
Copying test files...
Compiling...
Running...
----------------------------------------------------------
Test #1:
Input:
4
word
localization
internationalization
pneumonoultramicroscopicsilicovolcanoconiosis

Expected output:
word
l10n
i18n
p43s

Execution output:

==========================================================
 SUMMARY
==========================================================
  Test #1 [0.001s]: Wrong Answer
----------------------------------------------------------
 RESULT: Some tests are failing...
==========================================================
~$ 
```
naturally it will fail since we are not writing anything to the output. So, let's write [a simple solution](http://codeforces.com/contest/71/submission/19364315) 
for this problem and run `gocf test` again. This time, the test passed so you can submit already your work file as is.

Finally, you can archive your solution for historical purposes or for working on it later.
```
gocf archive
```

That's it!

Limitations
-----------
- tests cannot be run independently (i.e. run only test #3).
- the only supported judge by the import command is Codeforces.
- session properties cannot be changed, once created (you can change them manually, though).
- so far it works in Ubuntu 14.04 and OS X, using Go 1.6+. No idea if it works in other environments.
- memory limit is not taken into account. Suggestions on how to measure it would be appreciated.

Contributing
------------
Please, do contribute in any form you like: reporting bugs, suggesting features, improvements, ideas, docs, tests, ANYTHING. 
I would be extremely happy! :) Also, note that I started learning Go recently; as such, I could use a great deal of 
criticism and suggestions regarding the code.

License
-------
Licensed under Apache License 2.0 - http://www.apache.org/licenses/LICENSE-2.0
