every time you run it it gives an ever increasing number


$ go install github.com/jackdoe/updown/cmd/counter@latest



then you can use it for tail -f access.log|grep 123123 > tmp.$(counter).txt, so
you dont overwrite your temp work

