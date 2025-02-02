package charger

// Code generated by github.com/evcc-io/evcc/cmd/tools/decorate.go. DO NOT EDIT.

import (
	"github.com/evcc-io/evcc/api"
)

func decorateOpenWB20(base *OpenWB20, phaseSwitcher func(int) error, identifier func() (string, error)) api.Charger {
	switch {
	case identifier == nil && phaseSwitcher == nil:
		return base

	case identifier == nil && phaseSwitcher != nil:
		return &struct {
			*OpenWB20
			api.PhaseSwitcher
		}{
			OpenWB20: base,
			PhaseSwitcher: &decorateOpenWB20PhaseSwitcherImpl{
				phaseSwitcher: phaseSwitcher,
			},
		}

	case identifier != nil && phaseSwitcher == nil:
		return &struct {
			*OpenWB20
			api.Identifier
		}{
			OpenWB20: base,
			Identifier: &decorateOpenWB20IdentifierImpl{
				identifier: identifier,
			},
		}

	case identifier != nil && phaseSwitcher != nil:
		return &struct {
			*OpenWB20
			api.Identifier
			api.PhaseSwitcher
		}{
			OpenWB20: base,
			Identifier: &decorateOpenWB20IdentifierImpl{
				identifier: identifier,
			},
			PhaseSwitcher: &decorateOpenWB20PhaseSwitcherImpl{
				phaseSwitcher: phaseSwitcher,
			},
		}
	}

	return nil
}

type decorateOpenWB20IdentifierImpl struct {
	identifier func() (string, error)
}

func (impl *decorateOpenWB20IdentifierImpl) Identify() (string, error) {
	return impl.identifier()
}

type decorateOpenWB20PhaseSwitcherImpl struct {
	phaseSwitcher func(int) error
}

func (impl *decorateOpenWB20PhaseSwitcherImpl) Phases1p3p(p0 int) error {
	return impl.phaseSwitcher(p0)
}
