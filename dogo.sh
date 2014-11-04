if [ $# -ne 1 ]; then
    echo "valid number of arguments" $#
    echo "need number of arguments is "1
    exit 1
fi

cat << __EOT__
    argument is
    $1
    go install $1
__EOT__

go install $1

cat << __EOT__
    run $1

__EOT__

bin/$1

exit 0
