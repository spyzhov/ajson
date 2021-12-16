package jsonpath

import (
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

const __ = -1

// The action codes
const (
	cl internal.State = -2 /* colon           */
	cm internal.State = -3 /* comma           */
	qt internal.State = -4 /* quote           */
	bo internal.State = -5 /* bracket open    */
	co internal.State = -6 /* curly br. open  */
	bc internal.State = -7 /* bracket close   */
	cc internal.State = -8 /* curly br. close */
	ec internal.State = -9 /* curly br. empty */
	// jsonpath
	jd internal.State = -10 /* dollar */
	ja internal.State = -11 /* at */
	jq internal.State = -12 /* question mark */
	pc internal.State = -13 /* parentheses close */
	//	todo: function arguments!
	//	todo: implement ErrState and ErrLogic: `a + b * ()` - ErrLogic, `123)` - ErrState. Think about it...
)
