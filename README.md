#rx
---
## OverView
rx is rewrited command from gorm.

rx command is remove + backup command.

This command's main function are **remove + backup**, **undo rx**, **view trashbox**.

---
## Examples

- rx

  rx is rxmove + backup command in your trashbox. (Default: $HOME/.trashbox)

      $ ls a aa b
      a aa b
      $ rx a aa b
      $ ls
      $ ls $HOME/.trashbox
      a aa b

- rxls

  rxls is viewing item in your trashbox. Item name is patial mattch.

      $ rxls a
      --> ls a*
      (0) a original location : /home/kota999/a
      (1) aa original location : /home/kota999/aa

- rxundo

  rxundo is undo rx command. You select option in menu.

  Select menu is one option, not viewing item name.

      $ rxundo a
      --> recovering a*
      (0) a original location : /home/kota999/a
      (1) aa original location : /home/kota999/aa
      you will select one from (0, 1) > 1
      --> finish recovering to /home/kota999/aa
      $ ls aa
      aa

- rxbox

  rxbox is your trashbox operation. Your trashbox directory is writted in rx configure file. ($HOME/.rx)

  -box option: change your trashbox directory.

      $ ls
      $ cat $HOME/.rx
      /home/kota999/.trashbox
      $ rxbox -box $HOME/trash
      INFO: you setted trash-box directory option, so trash-box clear option is not effective.
      $ cat $HOME/.rx
      /home/kota999/trash
      $ ls
      trash

  -d option: initialize your trashbox configuration. ($HOME/.trashbox)

      $ ls
      trash
      $ cat $HOME/.rx
      /home/kota999/trash
      $ rxbox -d
      initializing trash-box configure and clear trash-box.
      INFO: initialized trash-box directory is $HOME/.trashbox.
      INFO: you setted trash-box directory option, so trash-box clear options is not effective.
      $ cat $HOME/.rx
      /home/kota999/.trashbox

  -c, -C option: clear all items in your trashbox directory.

      $ rxls
      --> ls *
      (0) a original location : /home/kota999/a
      (1) aa original location : /home/kota999/aa
      $ rxbox -C
      $ rxls
      --> ls *
      * is not match



---
## Commands / Options

    rx [-r|R|v|V] filename ...
    # remove + backup filename
    # options:
    #          -r or -R : remove + backup directory recursively
    #          -v or -V : view item name of interest

    rxundo [-v|V] filename ...
    # undo rx command
    # option:
    #          -v or -V : view detail (e.g. did backup date of filename)

    rxls [-l] filename ...
    # view items in your trashbox directory
    # option:
    #          -l : view detail (e.g. did backup date of filename)

    rxbox [-box|c|C|d] dirName ...
    # trashbox operation
    # options:
    #           -box : change your trashbox directory to dirName
    #           -c or -C : clear all items in your trashbox directory
    #           -d : initialize your trashbox configuration, trashbox directory is $HOME/.trashbox

---
## Installation


    git clone git://github.com/kota999/gorm.git

    go install rx rxbox rxls rxundo

You need export GOROOT, GOPATH in conjunction with the your environment.


---
## Infomation

* The config file of rx is $HOME/.rx, your trashbox path is written.
And the default trashbox directory is $HOME/.trashbox .

* Changing tarshbox directory manualy, you write your trashbox directory path in $HOME/.rx .

* In backup to your trash-box directory, rx command avoid fileName dupication . If there is a file test.txt and your command is

        $ rx test.txt

    , gorm command backup test.txt to test.txt.1 . And likewise, exist also test.txt.1 in your trash-box, backup test.txt to test.txt.2 .

* Remove + backup log of each items is writen in /your/trashbox/.prefix/filename.rx . (e.g. test.txt has test.txt.rx file)

* This command is not support the composition of multiple options, such as -rv.

