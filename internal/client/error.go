package client

import "fmt"

var ErrXidAlreadyTaken = fmt.Errorf("xid already taken")
var ErrInvalidClient = fmt.Errorf("client invalid")
