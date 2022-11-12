package cpanmeta

import (
	"encoding/json"
	"fmt"

	// local
	"github.com/cmburn/perlutils/internal"
)

// License indicates the license under which the distribution is released.
type License int

const (
	// LicenseUnknown is the default value for License, and implies that
	// the license info was not provided in the metadata.
	LicenseUnknown License = iota

	// LicenseOpenSource indicates an OSI-approved open source license not
	// covered by any of the other License constants.
	LicenseOpenSource

	// LicenseRestricted indicates a license requiring some form of
	// permission from the license holder to use the software.
	LicenseRestricted

	// LicenseUnrestricted indicates a license that does not require any
	// special permission, but is nonetheless not approved by the OSI.
	LicenseUnrestricted

	// LicenseAGPL3 indicates the software is licensed under the
	// GNU Affero General Public license, version 3.
	//
	// See also: https://opensource.org/licenses/agpl-3.0
	LicenseAGPL3

	// LicenseApache1_1 indicates the software is licensed under the
	// Apache license, version 1.1.
	//
	// See also: https://opensource.org/licenses/apache-1.1
	LicenseApache1_1

	// LicenseApache2_0 indicates the software is licensed under the
	// Apache license, version 2.0.
	//
	// See also: https://opensource.org/licenses/apache-2.0
	LicenseApache2_0

	// LicenseArtistic1 indicates the software is licensed under the
	// Artistic license, version 1.0.
	//
	// See also: https://opensource.org/licenses/artistic-1.0
	LicenseArtistic1

	// LicenseArtistic2 indicates the software is licensed under the
	// Artistic license, version 2.0.
	//
	// See also: https://opensource.org/licenses/artistic-2.0
	LicenseArtistic2

	// LicenseBSD indicates the software is licensed under the
	// 3-Clause BSD license.
	//
	// See also: https://opensource.org/licenses/BSD-3-Clause
	LicenseBSD

	// LicenseFreeBSD indicates the software is licensed under the
	// 2-Clause BSD license.
	//
	// See also: https://opensource.org/licenses/BSD-2-Clause
	LicenseFreeBSD

	// LicenseGFDL1_2 indicates the software is licensed under the
	// GNU Free Documentation license, version 1.2.
	//
	// See also: https://www.gnu.org/licenses/fdl-1.2.html
	LicenseGFDL1_2

	// LicenseGFDL1_3 indicates the software is licensed under the
	// GNU Free Documentation license, version 1.3.
	//
	// See also: https://www.gnu.org/licenses/fdl-1.3.html
	LicenseGFDL1_3

	// LicenseGPL1 indicates the software is licensed under the
	// GNU General Public license, version 1.
	//
	// See also: https://opensource.org/licenses/gpl-1.0
	LicenseGPL1

	// LicenseGPL2 indicates the software is licensed under the
	// GNU General Public license, version 2.
	//
	// See also: https://opensource.org/licenses/gpl-2.0
	LicenseGPL2

	// LicenseGPL3 indicates the software is licensed under the
	// GNU General Public license, version 3.
	//
	// See also: https://opensource.org/licenses/gpl-3.0
	LicenseGPL3

	// LicenseLGPL2_1 indicates the software is licensed under the
	// GNU Lesser General Public license, version 2.1.
	//
	// See also: https://opensource.org/licenses/LGPL-2.1
	LicenseLGPL2_1

	// LicenseLGPL3_0 indicates the software is licensed under the
	// GNU Lesser General Public license, version 3.0.
	//
	// See also: https://opensource.org/licenses/LGPL-3.0
	LicenseLGPL3_0

	// LicenseMIT indicates the software is licensed under the
	// MIT license.
	//
	// See also: https://opensource.org/licenses/MIT
	LicenseMIT

	// LicenseMozilla1_0 indicates the software is licensed under the
	// Mozilla Public license, version 1.0.
	//
	// See also: https://opensource.org/licenses/MPL-1.0
	LicenseMozilla1_0

	// LicenseMozilla1_1 indicates the software is licensed under the
	// Mozilla Public license, version 1.1.
	//
	// See also: https://opensource.org/licenses/MPL-1.1
	LicenseMozilla1_1

	// LicenseOpenSSL indicates the software is licensed under the
	// OpenSSL license.
	//
	// See also: https://www.openssl.org/source/license-openssl-ssleay.txt
	LicenseOpenSSL

	// LicensePerl5 indicates the software is licensed under the
	// [same terms as Perl 5 itself]. That is, LicensePerl5 is equivalent
	// to LicenseGPL1 and LicenseArtistic1.
	//
	// [same terms as Perl 5 itself]: https://dev.perl.org/licenses/
	LicensePerl5

	// LicenseQPL1_0 indicates the software is licensed under the
	// Q Public license, version 1.0.
	//
	// See also: https://opensource.org/licenses/QPL-1.0
	LicenseQPL1_0

	// LicenseSSLeay indicates the software is licensed under the
	// SSLeay license.
	//
	// See also: https://www.openssl.org/source/license-openssl-ssleay.txt
	LicenseSSLeay

	// LicenseSun indicates the software is licensed under the
	// Sun Industry Standards Source license.
	//
	// See also: https://opensource.org/licenses/SISSL
	LicenseSun

	// LicenseZlib indicates the software is licensed under the
	// zlib license.
	//
	// See also: https://opensource.org/licenses/Zlib
	LicenseZlib
)

func (l *License) String() string {
	switch *l {
	case LicenseOpenSource:
		return "open_source"
	case LicenseRestricted:
		return "restricted"
	case LicenseUnrestricted:
		return "unrestricted"
	case LicenseAGPL3:
		return "agpl_3"
	case LicenseApache1_1:
		return "apache_1_1"
	case LicenseApache2_0:
		return "apache_2_0"
	case LicenseArtistic1:
		return "artistic_1"
	case LicenseArtistic2:
		return "artistic_2"
	case LicenseBSD:
		return "bsd"
	case LicenseFreeBSD:
		return "freebsd"
	case LicenseGFDL1_2:
		return "gfdl_1_2"
	case LicenseGFDL1_3:
		return "gfdl_1_3"
	case LicenseGPL1:
		return "gpl_1"
	case LicenseGPL2:
		return "gpl_2"
	case LicenseGPL3:
		return "gpl_3"
	case LicenseLGPL2_1:
		return "lgpl_2_1"
	case LicenseLGPL3_0:
		return "lgpl_3_0"
	case LicenseMIT:
		return "mit"
	case LicenseMozilla1_0:
		return "mozilla_1_0"
	case LicenseMozilla1_1:
		return "mozilla_1_1"
	case LicenseOpenSSL:
		return "openssl"
	case LicensePerl5:
		return "perl_5"
	case LicenseQPL1_0:
		return "qpl_1_0"
	case LicenseSSLeay:
		return "ssleay"
	case LicenseSun:
		return "sun"
	case LicenseZlib:
		return "zlib"
	case LicenseUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

func NewLicense(str string) (License, error) {
	switch str {
	case "unknown":
		return LicenseUnknown, nil
	case "open_source":
		return LicenseOpenSource, nil
	case "restricted":
		return LicenseRestricted, nil
	case "unrestricted":
		return LicenseUnrestricted, nil
	case "agpl_3":
		return LicenseAGPL3, nil
	case "apache_1_1":
		return LicenseApache1_1, nil
	case "apache_2_0":
		return LicenseApache2_0, nil
	case "artistic_1":
		return LicenseArtistic1, nil
	case "artistic_2":
		return LicenseArtistic2, nil
	case "bsd":
		return LicenseBSD, nil
	case "freebsd":
		return LicenseFreeBSD, nil
	case "gfdl_1_2":
		return LicenseGFDL1_2, nil
	case "gfdl_1_3":
		return LicenseGFDL1_3, nil
	case "gpl_1":
		return LicenseGPL1, nil
	case "gpl_2":
		return LicenseGPL2, nil
	case "gpl_3":
		return LicenseGPL3, nil
	case "lgpl_2_1":
		return LicenseLGPL2_1, nil
	case "lgpl_3_0":
		return LicenseLGPL3_0, nil
	case "mit":
		return LicenseMIT, nil
	case "mozilla_1_0":
		return LicenseMozilla1_0, nil
	case "mozilla_1_1":
		return LicenseMozilla1_1, nil
	case "openssl":
		return LicenseOpenSSL, nil
	case "perl_5":
		return LicensePerl5, nil
	case "qpl_1_0":
		return LicenseQPL1_0, nil
	case "ssleay":
		return LicenseSSLeay, nil
	case "sun":
		return LicenseSun, nil
	case "zlib":
		return LicenseZlib, nil
	default:
		return LicenseUnknown, fmt.Errorf("unknown license: %s", str)
	}
}

func (l *License) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	license, err := NewLicense(str)
	if err != nil {
		return err
	}
	*l = license
	return nil
}

func (l *License) MarshalJSON() ([]byte, error) {
	return internal.WrapEnumTypeJSON(l)
}
