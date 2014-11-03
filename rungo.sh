if [ $# -ne 1 ]; then
    echo "valid number of arguments" $#
    echo "need number of arguments is "1
    exit 1
fi

cat << __EOT__
    argument is
    $1
    Run $1.go

__EOT__

go run src/github.com/user/$1/$1.go

exit 0
