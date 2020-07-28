package common

import (
    "errors"
    "testing"
)

func TestPrintWrappedErrors(t *testing.T) {
    rootErr := NewHttpError(nil, "web/server/common", "This is a test (root)", 0)
    goErr := errors.New("dummy go error")
    wrapGoErr := NewHttpError(goErr, "web/server/common", "This is another test (wrap go)", 1)
    wrapHttpErr := NewHttpError(rootErr, "web/server/common", "This is the final (?) test (wrap http)", 2)
    doubleWrapHttpErr := NewHttpError(wrapHttpErr, "web/server/common", "Final test (for real)", 3)

    t.Logf("root err: %+v\n", rootErr)
    t.Logf("wrap go err: %+v\n", wrapGoErr)
    t.Logf("wrap http err: %+v\n", wrapHttpErr)
    t.Logf("double wrap http err: %+v\n", doubleWrapHttpErr)
}
