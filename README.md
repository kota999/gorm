gorm command is remove + backup command.
Rename thrw ---> gorm .


Simple Useage :


    gorm [fileName]       ... remove filename and backup in $HOME/.trashbox (default)

    gorm -r [fileName]    ... remove directory and backup in $HOME/.trashbox (default)

    gorm -box=[dirName]   ... register your trash-box directory

    gorm -c               ... clear your trash-box


Installation :


    git clone git://github.com/kota999/gorm.git

    go install gorm

You need export GOROOT, GOPATH in conjunction with the your environment.


Options :


  * For Remove options

        -r, -R              ... remove and backup all files and directories in dirName recursively

        -v                  ... print names of remove and backup all files and directories


  * For backup options

        -box                ... register your trash-box directory
                                if use this option, gorm command skip other options

        -c, -C              ... clear your trash-box directory



Infomation :

* The config file of gorm is $HOME/.gorm, your trash-box path is written by .gorm .
And the default trash-box directory is $HOME/.trashbox .

* In backup to your trash-box directory, gorm command avoid fileName dupication . If there is a file test.txt and your command is

        gorm test.txt

    , gorm command backup test.txt to test.txt.1 . And likewise, exist also test.txt.1 in your trash-box, backup test.txt to test.txt.2 .

* This command is not support the composition of multiple options, such as -rv.

