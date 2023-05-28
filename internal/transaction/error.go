package transaction

import "fmt"

var ErrReferenceNoAlreadyExists = fmt.Errorf("reference number already exists")
var ErrEmptyCustomerXid = fmt.Errorf("customer xid is required")
var ErrEmptyReferenceId = fmt.Errorf("reference id is required")
