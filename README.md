rt command is remove + backup command.
Rename thrw ---> rt .


Simple Useage :


    rt [fileName]       ... remove filename and backup in $HOME/.trashbox (default)

    rt -r [fileName]    ... remove directory and backup in $HOME/.trashbox (default)

    rt -box=[dirName]   ... register your trash-box directory

    rt -c               ... clear your trash-box


Installation :


    git clone git://github.com/kota999/rt.git

    go install rt

You need export GOROOT, GOPATH in conjunction with the your environment.


Options :


  * For Remove options

        -r, -R              ... remove and backup all files and directories in dirName recursively

        -v                  ... print names of remove and backup all files and directories


  * For backup options

        -box                ... register your trash-box directory
                                if use this option, rt command skip other options

        -c, -C              ... clear your trash-box directory



Infomation :

* The config file of rt is $HOME/.thrw, your trash-box path is written by .thrw .
And the default trash-box directory is $HOME/.trashbox .

* In backup to your trash-box directory, rt command avoid fileName dupication . If there is a file test.txt and your command is

        rt test.txt

    , rt command backup test.txt to test.txt.1 . And likewise, exist also test.txt.1 in your trash-box, backup test.txt to test.txt.2 .

* This command is not support the synthesis of multiple options, such as -rv.

