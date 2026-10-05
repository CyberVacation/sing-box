//go:build !linux

package main

import "github.com/CyberVacation/rostra/option"

func runInUserNamespaceIfNeeded(options option.Options, optionsList []*OptionsEntry) error {
	return nil
}
