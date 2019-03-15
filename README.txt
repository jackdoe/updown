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




