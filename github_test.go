package github

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBuildPath(t *testing.T) {
	Convey("buildPath works", t, func() {
		So(buildPath("hello", "joe"), ShouldEqual, "hello/joe")
		So(buildPath("hello/", "joe"), ShouldEqual, "hello/joe")
		So(buildPath("hello", "/joe"), ShouldEqual, "hello/joe")
		So(buildPath("hello/", "/joe"), ShouldEqual, "hello/joe")
		So(buildPath("/hello/", "/joe"), ShouldEqual, "/hello/joe")
	})
}
