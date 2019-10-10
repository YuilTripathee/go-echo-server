package main

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

const forceDisconnectAfter = time.Second * 5

var (
	// ErrBadFormat is error format introduced for the mail having invalid format
	ErrBadFormat = errors.New("invalid format")
	// ErrUnresolvableHost is error format introduced for incorrect (dummy) email addresses
	ErrUnresolvableHost = errors.New("unresolvable host")
	// ErrNumNotInRange is error format for given number not in range
	ErrNumNotInRange = errors.New("number out of range")

	emailRegExp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// SMTPError is defined to structure new SMTP error formats
type SMTPError struct {
	Err error
}

// Error() returns error message for new SMTPError
func (e SMTPError) Error() string {
	return e.Err.Error()
}

// Code returns error code for new SMTPError
func (e SMTPError) Code() string {
	return e.Err.Error()[0:3]
}

// NewError returns new error with error message and error code
func NewError(err error) SMTPError {
	return SMTPError{
		Err: err,
	}
}

// isEmail checks if provided string is a valid email address and returns errors if present
func isEmail(TestString string) error {
	if !emailRegExp.MatchString(TestString) {
		return ErrBadFormat
	}
	return nil
}

// isValidatedHost checks if provided email string is a valid email address
func isValidatedHost(TestEmailString string) error {
	_, host := splitMailAddr(TestEmailString)
	mx, err := net.LookupMX(host)
	if err != nil {
		return ErrUnresolvableHost
	}

	client, err := DialTimeout(fmt.Sprintf("%s:%d", mx[0].Host, 25), forceDisconnectAfter)
	if err != nil {
		return NewError(err)
	}
	defer client.Close()

	err = client.Rcpt(TestEmailString)
	if err != nil {
		return NewError(err)
	}
	return nil
}

// DialTimeout returns a new Client connected to an SMTP server at addr.
// The addr must include a port, as in "mail.example.com:smtp".
func DialTimeout(addr string, timeout time.Duration) (*smtp.Client, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	t := time.AfterFunc(timeout, func() { conn.Close() })
	defer t.Stop()

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// splitMailAddr returns account and host string from given email address
func splitMailAddr(email string) (account, host string) {
	i := strings.LastIndexByte(email, '@')
	account = email[:i]
	host = email[i+1:]
	return
}

func isNumberInRange(minValue int, maxValue int, testNumber int) error {
	if testNumber < minValue || testNumber > maxValue {
		return ErrNumNotInRange
	}
	return nil
}

// for advanced operations check package: govalidator
// this is much proven way to implement validation in Golang
