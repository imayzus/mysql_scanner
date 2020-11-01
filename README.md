# mysql_scanner

This project provides a script, scan_for_mysql.go, to check if mysql is running on a particular IP and port.

It will retrieve a banner from the port and check if it contains 'mysql'. It will also run nmap (if it's installed on the system) and check if its
output for service is 'mysql'.

The program has been tested on osx with mysql running on port 3306 (the default port).

To check if mysql is running on a particular port, run:
go run scan_for_mysql.go <hostName> <port>

For example, to check if the localhost has mysql runnning on port 3306, run:
go run scan_for_mysql.go 127.0.0.1 3306

if parameters are omitted, the program will run with default values, e.g.:
go run scan_for_mysql.go

To run tests, use "go test"
(when in the same directory as scan_for_mysql.go and scan_for_mysql_test.go)
