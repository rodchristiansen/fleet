package oval_parsed

import (
	"fmt"
	"strings"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

type ObjectInfoState struct {
	Name           *ObjectStateString      `json:",omitempty"`
	Arch           *ObjectStateString      `json:",omitempty"`
	Epoch          *ObjectStateSimpleValue `json:",omitempty"`
	Release        *ObjectStateSimpleValue `json:",omitempty"`
	Version        *ObjectStateSimpleValue `json:",omitempty"`
	Evr            *ObjectStateEvrString   `json:",omitempty"`
	SignatureKeyId *ObjectStateString      `json:",omitempty"`
	ExtendedName   *ObjectStateString      `json:",omitempty"`
	FilePath       *ObjectStateString      `json:",omitempty"`
	Operator       OperatorType            `json:"operator"`
}

// EvalSoftware evaluates the passed software against the specified state.
func (sta ObjectInfoState) EvalSoftware(s fleet.Software) (bool, error) {
	var results []bool

	if sta.Name != nil {
		rEval, err := sta.Name.Eval(s.Name)
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.Arch != nil {
		rEval, err := sta.Arch.Eval(s.Arch)
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.Epoch != nil {
		rEval, err := sta.Epoch.Eval(fmt.Sprint(epoch(s.Version)))
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.Release != nil {
		rEval, err := sta.Release.Eval(release(s.Version))
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.Version != nil {
		rEval, err := sta.Version.Eval(version(s.Version))
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.Evr != nil {
		rEval, err := sta.Evr.Eval(s.Version, Rpmvercmp)
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.SignatureKeyId != nil {
		// Assume that all installed software was signed by the proper third party (RedHat), we are
		// doing this basically because there's no way to get the signature key ATM and even if we
		// have it we want to reuse the RHEL OVAL definitions for CentOS
		results = append(results, true)
	}

	if len(results) == 0 {
		return false, fmt.Errorf("invalid empty state")
	}

	return sta.Operator.Eval(results...), nil
}

func (sta ObjectInfoState) EvalOSVersion(version fleet.OSVersion) (bool, error) {
	var results []bool

	// If 'sta' is used for specifying the state of a RpmVerifyFile test, 'Name' refers to the name of the
	// file, when making assertions against the installed OS, the file in question will be
	// /etc/redhat-release, so in order to use the same test for CentOS distros, we will need to
	// normalize the value.
	if sta.Name != nil {
		var nName string
		if version.Platform == "rhel" {
			nName = "redhat-release"
		}
		rEval, err := sta.Name.Eval(nName)
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if sta.Version != nil {
		pName := strings.Trim(version.Name, " ")
		pVer := pName[strings.LastIndex(pName, " ")+1:]
		rEval, err := sta.Version.Eval(pVer)
		if err != nil {
			return false, err
		}
		results = append(results, rEval)
	}

	if len(results) == 0 {
		return false, fmt.Errorf("invalid empty state")
	}

	return sta.Operator.Eval(results...), nil
}
