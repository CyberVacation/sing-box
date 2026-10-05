//go:build !with_utls

package v2rayxhttp

import "github.com/CyberVacation/rostra/common/tls"

func isReality(config tls.Config) bool { return false }
