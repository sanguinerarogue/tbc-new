package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func (warlock *Warlock) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_WarlockAssignedCurseIsActive:
		return warlock.newValueWarlockAssignedCurseIsActive(rot, config.GetWarlockAssignedCurseIsActive())
	default:
		return nil
	}
}

type APLValueWarlockAssignedCurseIsActive struct {
	core.DefaultAPLValueImpl
	warlock *Warlock
}

func (x *APLValueWarlockAssignedCurseIsActive) GetInnerActions() []*core.APLAction { return nil }
func (x *APLValueWarlockAssignedCurseIsActive) GetAPLValues() []core.APLValue      { return nil }
func (x *APLValueWarlockAssignedCurseIsActive) Finalize(*core.APLRotation)         {}
func (x *APLValueWarlockAssignedCurseIsActive) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}
func (x *APLValueWarlockAssignedCurseIsActive) GetSpellFromAction(sim *core.Simulation) *core.Spell {
	switch x.warlock.Options.CurseOptions {
	case proto.WarlockOptions_Agony:
		return x.warlock.CurseOfAgony

	case proto.WarlockOptions_Doom:
		if sim.BaseDuration-sim.CurrentTime > time.Minute*1 {
			return x.warlock.CurseOfDoom
		} else {
			return x.warlock.CurseOfAgony
		}

	case proto.WarlockOptions_Elements:
		return x.warlock.CurseOfElements

	case proto.WarlockOptions_Recklessness:
		return x.warlock.CurseOfRecklessness
	}

	return nil
}

// func (x APLValueWarlockAssignedCurse)

func (warlock *Warlock) newValueWarlockAssignedCurseIsActive(rot *core.APLRotation, config *proto.APLValueWarlockAssignedCurseIsActive) core.APLValue {
	return &APLValueWarlockAssignedCurseIsActive{
		warlock: warlock,
	}
}

func (x *APLValueWarlockAssignedCurseIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}

func (x *APLValueWarlockAssignedCurseIsActive) GetBool(sim *core.Simulation) bool {
	assignedCurse := x.GetSpellFromAction(sim)
	aura := x.warlock.CurrentTarget.GetAuraByID(assignedCurse.ActionID)

	return aura.IsActive()
}

func (x *APLValueWarlockAssignedCurseIsActive) String() string {
	return "Cast Assigned Curse"
}

func (warlock *Warlock) NewAPLAction(rot *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CastWarlockAssignedCurse:
		return warlock.newActionWarlockAssignedCurseAction(rot, config.GetCastWarlockAssignedCurse())
	default:
		return nil
	}
}

type APLActionCastWarlockAssignedCurse struct {
	warlock    *Warlock
	lastAction time.Duration
}

func (x *APLActionCastWarlockAssignedCurse) GetInnerActions() []*core.APLAction { return nil }
func (x *APLActionCastWarlockAssignedCurse) GetAPLValues() []core.APLValue      { return nil }
func (x *APLActionCastWarlockAssignedCurse) Finalize(*core.APLRotation)         {}
func (x *APLActionCastWarlockAssignedCurse) PostFinalize(*core.APLRotation)     {}
func (x *APLActionCastWarlockAssignedCurse) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}
func (x *APLActionCastWarlockAssignedCurse) ReResolveVariableRefs(*core.APLRotation, map[string]*proto.APLValue) {
}

func (x *APLActionCastWarlockAssignedCurse) GetSpellFromAction(sim *core.Simulation) *core.Spell {

	switch x.warlock.Options.CurseOptions {
	case proto.WarlockOptions_Agony:
		return x.warlock.CurseOfAgony

	case proto.WarlockOptions_Doom:
		if sim.BaseDuration-sim.CurrentTime > time.Minute*1 {
			return x.warlock.CurseOfDoom
		} else {
			return x.warlock.CurseOfAgony
		}

	case proto.WarlockOptions_Elements:
		return x.warlock.CurseOfElements

	case proto.WarlockOptions_Recklessness:
		return x.warlock.CurseOfRecklessness
	}

	return nil
}

func (warlock *Warlock) newActionWarlockAssignedCurseAction(_ *core.APLRotation, _ *proto.APLActionCastWarlockAssignedCurse) core.APLActionImpl {
	return &APLActionCastWarlockAssignedCurse{
		warlock: warlock,
	}
}

func (x *APLActionCastWarlockAssignedCurse) Execute(sim *core.Simulation) {
	x.GetSpellFromAction(sim).Cast(sim, x.warlock.CurrentTarget)
}

func (x *APLActionCastWarlockAssignedCurse) IsReady(sim *core.Simulation) bool {
	return x.GetSpellFromAction(sim).CanCast(sim, x.warlock.CurrentTarget)
}

func (x *APLActionCastWarlockAssignedCurse) Reset(*core.Simulation) {
	x.lastAction = core.DurationFromSeconds(-100)
}

func (x *APLActionCastWarlockAssignedCurse) String() string {
	return "Cast Assigned Curse"
}
