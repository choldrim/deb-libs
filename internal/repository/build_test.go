package packages

import "testing"
import "bytes"
import C "gopkg.in/check.v1"

type testWrap struct{}

func Test(t *testing.T) { C.TestingT(t) }
func init() {
	C.Suite(&testWrap{})
}

func (*testWrap) TestParsingDBComponent(c *C.C) {
	types, err := parsePackageDBComponent("testdata/Packages")
	c.Check(err, C.Equals, nil)
	c.Check(len(types), C.Equals, 523)
}

func (*testWrap) TestBuildCache(c *C.C) {
	targetDir := "testdata/packages"
	codeName := "unstable"

	m, err := NewManager("testdata", "http://pools.corp.deepin.com/deepin/", codeName)
	c.Check(err, C.Equals, nil)
	err = m.UpdateDB()
	c.Check(err, C.Equals, nil)

	t, ok := m.Get("lastore-daemon")
	c.Check(ok, C.Equals, true)
	c.Check(t.Homepage, C.Equals, "http://github.com/linuxdeepin/lastore-daemon")

	rf, err := GetReleaseFile(targetDir, codeName)
	c.Check(err, C.Equals, nil)
	err = BuildCache(rf, targetDir, targetDir)
	c.Check(err, C.Equals, nil)

}

var d = `Package: aac-enc
Source: fdk-aac
Version: 0.1.3+20140816-2
Installed-Size: 705
Maintainer: Debian Multimedia Maintainers <pkg-multimedia-maintainers@lists.alioth.debian.org>
Architecture: amd64
Depends: libfdk-aac0 (= 0.1.3+20140816-2), libc6 (>= 2.4)
Size: 666554
SHA256: d09f8c35f8817bc67b67ebc7af94d7b26ba656af2bea4ed579e13e03db718cee
SHA1: b9a70c3b65f7ad6b62f56c2b8cc916b156c38713
MD5sum: 9703f7d0d4463b198bfd57b45fefd8ab
Description: Fraunhofer FDK AAC Codec Library - frontend binary
 test multiline
Homepage: https://github.com/mstorsjo/fdk-aac
Description-md5: 16f812d0c8b3e09448f6f7d88536e135
Section: non-free/sound
Priority: optional
Filename: pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb
`

func (*testWrap) TestType(c *C.C) {
	b := bytes.NewBufferString(d)
	t, err := buildType(b)
	c.Assert(err, C.Equals, nil)
	c.Check(t.Filename, C.Equals, "pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb")
	c.Check(t.Size, C.Equals, 666554)
	//c.Check(t.Description, C.Equals, "Description: Fraunhofer FDK AAC Codec Library - frontend binary test multiline")
}

func (*testWrap) TestDSC(c *C.C) {

	b := bytes.NewBufferString(d)
	d, err := NewDSCFile(b)
	c.Check(err, C.Equals, nil)

	c.Check(d.GetString("Package"), C.Equals, "aac-enc")

	c.Check(d.GetString("Source"), C.Equals, "fdk-aac")

	c.Check(d.GetString("Version"), C.Equals, "0.1.3+20140816-2")

	c.Check(d.GetString("installed-size"), C.Equals, "705")

	c.Check(d.GetString("archiTecTure"), C.Equals, "amd64")

	c.Check(d.GetString("depends"), C.Equals, "libfdk-aac0 (= 0.1.3+20140816-2), libc6 (>= 2.4)")

	c.Check(d.GetString("description"), C.Equals, `Fraunhofer FDK AAC Codec Library - frontend binary
 test multiline`)

	c.Check(d.GetString("priority"), C.Equals, "optional")

	c.Check(d.GetString("Filename"), C.Equals, "pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb")
}
