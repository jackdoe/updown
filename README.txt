go get github.com/jackdoe/updown

go install github.com/jackdoe/updown/cmd/snaketocamel
go install github.com/jackdoe/updown/cmd/cameltosnake

% ls -1 ~/work/baxx/help/t
email_after_registration.txt
email_notification.txt
email_payment_cancel.txt
email_payment_thanks.txt
% ls -1 ~/work/baxx/help/t | snaketocamel
EmailAfterRegistrationtxt
EmailNotificationtxt
EmailPaymentCanceltxt
EmailPaymentThankstxt
EmailValidationtxt

-------

go install github.com/jackdoe/updown/cmd/encrypt

use https://github.com/minio/sio to encrypt streams

encrypt -h
Usage of ./encrypt:
  -d    decrypt
  -k string
        encryption key, filename or '-' to ask; key is sha256(value) (default "-")


% echo 'hello world' | encrypt -k .pass | encrypt -d -k .pass
hello world
% echo 'hello world' | encrypt > encrypted
Key:
% cat encrypted | encrypt -d
Key:
hello world

% cat /dev/zero | encrypt -k .pass | encrypt -d -k .pass | pv > /dev/null
25.8GiB 0:00:10 [1.16GiB/s] [            <=>                                                  ]

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND                  
23296 jack      20   0   10156   6240   1484 S  80.0   0.0   0:34.21 encrypt                  
23297 jack      20   0   10156   6276   1540 S  80.0   0.0   0:33.62 encrypt    

-------

go install github.com/jackdoe/updown/cmd/plot
cat > example
1
2
4
5
^D
cat example | plot > example.png


-------
go install github.com/jackdoe/updown/cmd/sendme

%  sendme -html a.html -text a.txt -u user@gmail.com -subject test -to user2@gmail.com,user3@gmail.com

------
go install github.com/jackdoe/updown/cmd/pj

% echo '{"a":1}' | pj
{
  "a": 1
}


------
go install github.com/jackdoe/updown/cmd/py

% cat <<EOL | py
apiVersion: v1
metadata:
       name: example
EOL

apiVersion: v1
metadata:
  name: example


------
go install github.com/jackdoe/updown/cmd/ph

% echo '<html><body><div>hello</div></body></html>' | ph
<html>
  <body>
    <div>
      hello
    </div>
  </body>
</html


------
go install github.com/jackdoe/updown/cmd/onc

for each connection execute the command and pass the connection output
as the command input, good example is executing tar for each
connection

% onc -l 7000 tar xf
...

% tar -cf - . | nc -N localhost 7000
% tar -cf - . | nc -N localhost 7000
% tar -cf - . | nc -N localhost 7000


------
go install github.com/jackdoe/updown/cmd/quant

quantiles from stdin

% cat <<EOF | quant
1
2
4
2
4
3
2
EOF
perc25: 1.000000
perc50: 2.000000
perc75: 3.000000
perc90: 4.000000
perc99: 4.000000
count: 7

------
go install github.com/jackdoe/updown/cmd/sumint

sum integers

% cat <<EOF | sumint
1
2
4
2
4
3
2
EOF

16
%

-----

go install github.com/jackdoe/updown/cmd/pagerank

cat <<EOF | pagerank
a b
a c
c b
b d
EOF
0.1205 a
0.3176 b
0.1716 c
0.3903 d

----
go install github.com/jackdoe/updown/cmd/groupby

group csv by summing all coumns based on column
e.g.

cat <<EOF | groupby -separator ';'
a;1;2;3
b;1;2;3
a;2;3;4
EOF
a;3.00;5.00;7.00
b;1.00;2.00;3.00

---

----
go install github.com/jackdoe/updown/cmd/delta

group csv by summing all coumns based on column
e.g.

cat <<EOF | delta
1
2
3
1
5
10
1000
EOF

1
1
1
-2
4
5
990

---