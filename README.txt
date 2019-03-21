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
