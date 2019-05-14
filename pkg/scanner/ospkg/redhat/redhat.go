package redhat

import (
	"strings"

	"github.com/knqyf263/fanal/analyzer"
	version "github.com/knqyf263/go-rpm-version"
	"github.com/knqyf263/trivy/pkg/log"
	"github.com/knqyf263/trivy/pkg/scanner/utils"
	"github.com/knqyf263/trivy/pkg/vulnsrc/redhat"
	"github.com/knqyf263/trivy/pkg/vulnsrc/vulnerability"
	"golang.org/x/xerrors"
)

type Scanner struct{}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Detect(osVer string, pkgs []analyzer.Package) ([]vulnerability.DetectedVulnerability, error) {
	log.Logger.Info("Detecting RHEL/CentOS vulnerabilities...")
	if strings.Count(osVer, ".") > 0 {
		osVer = osVer[:strings.Index(osVer, ".")]
	}
	log.Logger.Debugf("redhat: os version: %s", osVer)
	log.Logger.Debugf("redhat: the number of packages: %s", len(pkgs))

	var vulns []vulnerability.DetectedVulnerability
	for _, pkg := range pkgs {
		advisories, err := redhat.Get(osVer, pkg.SrcName)
		if err != nil {
			return nil, xerrors.Errorf("failed to get Red Hat advisories: %w", err)
		}

		installed := utils.FormatSrcVersion(pkg)
		installedVersion := version.NewVersion(installed)
		for _, adv := range advisories {
			fixedVersion := version.NewVersion(adv.FixedVersion)

			if installedVersion.LessThan(fixedVersion) {
				vuln := vulnerability.DetectedVulnerability{
					VulnerabilityID:  adv.VulnerabilityID,
					PkgName:          pkg.Name,
					InstalledVersion: installed,
					FixedVersion:     adv.FixedVersion,
				}
				vulns = append(vulns, vuln)
			}
		}
	}
	return vulns, nil
}
