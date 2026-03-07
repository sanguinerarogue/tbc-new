package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const (
	HunterBaseMaxRange         = 35
	ThoridalTheStarsFuryItemID = 34334
)

var TalentTreeSizes = [3]int{21, 20, 24}

type Hunter struct {
	core.Character

	ClassSpellScaling float64

	Talents *proto.HunterTalents
	Options *proto.HunterOptions

	Pet *HunterPet

	AmmoDPS         float64
	AmmoDamageBonus float64

	killCommandEnabledUntil time.Duration // Time that KC enablement expires.

	AimedShot        *core.Spell
	ArcaneShot       *core.Spell
	AspectOfTheHawk  *core.Spell
	AspectOfTheViper *core.Spell
	BestialWrath     *core.Spell
	KillCommand      *core.Spell
	MultiShot        *core.Spell
	RapidFire        *core.Spell
	RaptorStrike     *core.Spell
	Readiness        *core.Spell
	ScorpidSting     *core.Spell
	SerpentSting     *core.Spell
	SteadyShot       *core.Spell
	// HuntersMarkSpell *core.Spell

	AspectOfTheHawkAura  *core.Aura
	AspectOfTheViperAura *core.Aura
	GronnStalker2PcAura  *core.Aura
	RapidFireAura        *core.Aura
	TalonOfAlarAura      *core.Aura
	TheBeastWithinAura   *core.Aura
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func RegisterHunter() {
	core.RegisterAgentFactory(
		proto.Player_Hunter{},
		proto.Spec_SpecHunter,
		func(character *core.Character, options *proto.Player, raid *proto.Raid) core.Agent {
			return NewHunter(character, options, options.GetHunter().Options.ClassOptions, raid)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Hunter)
			if !ok {
				panic("Invalid spec value for Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewHunter(character *core.Character, options *proto.Player, hunterOptions *proto.HunterOptions, raid *proto.Raid) *Hunter {
	hunter := &Hunter{
		Character: *character,
		Talents:   &proto.HunterTalents{},
		Options:   hunterOptions,
	}

	if hunter.Options.PetType == proto.HunterOptions_Bat || hunter.Options.PetType == proto.HunterOptions_Owl {
		raid.Debuffs.Screech = false
	}

	if hunter.Talents.ExposeWeakness > 0 {
		raid.Debuffs.ExposeWeaknessHunterAgility = 0
		raid.Debuffs.ExposeWeaknessUptime = 0
	}

	core.FillTalentsProto(hunter.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	hunter.PseudoStats.CanParry = true

	hunter.EnableManaBar()

	rangedWeapon := hunter.WeaponFromRanged(hunter.DefaultMeleeCritMultiplier())
	hunter.PseudoStats.RangedSpeedMultiplier = 1.0
	if wep := hunter.GetRangedWeapon(); wep != nil && wep.ID == ThoridalTheStarsFuryItemID {
		hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
	} else {
		switch hunter.Options.Ammo {
		case proto.HunterOptions_TimelessArrow:
			hunter.AmmoDPS = 53
		case proto.HunterOptions_MysteriousArrow:
			hunter.AmmoDPS = 46.5
		case proto.HunterOptions_AdamantiteStinger:
			hunter.AmmoDPS = 43
		case proto.HunterOptions_WardensArrow:
			hunter.AmmoDPS = 37
		case proto.HunterOptions_HalaaniRazorshaft:
			hunter.AmmoDPS = 34
		case proto.HunterOptions_BlackflightArrow:
			hunter.AmmoDPS = 32
		}
		hunter.AmmoDamageBonus = hunter.AmmoDPS * rangedWeapon.SwingSpeed
		rangedWeapon.BaseDamageMin += hunter.AmmoDamageBonus
		rangedWeapon.BaseDamageMax += hunter.AmmoDamageBonus

		switch hunter.Options.QuiverBonus {
		case proto.HunterOptions_Speed10:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.1
		case proto.HunterOptions_Speed11:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.11
		case proto.HunterOptions_Speed12:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.12
		case proto.HunterOptions_Speed13:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.13
		case proto.HunterOptions_Speed14:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.14
		case proto.HunterOptions_Speed15:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
		}
	}

	hunter.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotRanged}, func(sim *core.Simulation, slot proto.ItemSlot) {
		ranged := hunter.AutoAttacks.Ranged()
		if ranged == nil {
			return
		}

		// Remove old ammo bonus
		ranged.BaseDamageMin -= hunter.AmmoDamageBonus
		ranged.BaseDamageMax -= hunter.AmmoDamageBonus

		if hunter.GetRangedWeapon() != nil && hunter.GetRangedWeapon().ID == ThoridalTheStarsFuryItemID {
			hunter.PseudoStats.RangedSpeedMultiplier = 1.15
			hunter.AmmoDamageBonus = 0
		} else {
			hunter.AmmoDamageBonus = hunter.AmmoDPS * ranged.SwingSpeed
			ranged.BaseDamageMin += hunter.AmmoDamageBonus
			ranged.BaseDamageMax += hunter.AmmoDamageBonus
		}
	})

	hunter.EnableAutoAttacks(hunter, core.AutoAttackOptions{
		Ranged:          rangedWeapon,
		MainHand:        hunter.WeaponFromMainHand(hunter.DefaultMeleeCritMultiplier()),
		OffHand:         hunter.WeaponFromOffHand(hunter.DefaultMeleeCritMultiplier()),
		ReplaceMHSwing:  hunter.TryRaptorStrike,
		AutoSwingRanged: true,
		AutoSwingMelee:  true,
	})

	rangedConfig := hunter.AutoAttacks.RangedConfig()
	rangedConfig.MaxRange = HunterBaseMaxRange

	hunter.AddStatDependencies()

	hunter.Pet = hunter.NewHunterPet()

	return hunter
}

func (hunter *Hunter) RegisterRangedSpell(config core.SpellConfig) *core.Spell {
	if config.Cast.ModifyCast == nil {
		config.Cast.ModifyCast = func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
			cast.CastTime = spell.CastTime()
			hunter.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime+cast.CastTime)
		}
	}

	if config.Cast.CastTime == nil {
		config.Cast.CastTime = func(spell *core.Spell) time.Duration {
			return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.TotalRangedHasteMultiplier())
		}
	}

	return hunter.RegisterSpell(config)
}

func (hunter *Hunter) Initialize() {
	hunter.AutoAttacks.MHConfig().CritMultiplier = hunter.DefaultMeleeCritMultiplier()
	hunter.AutoAttacks.OHConfig().CritMultiplier = hunter.DefaultMeleeCritMultiplier()
	hunter.AutoAttacks.RangedConfig().CritMultiplier = hunter.DefaultMeleeCritMultiplier()

	hunter.RegisterSpells()
	hunter.addPvpGloves()
}

func (hunter *Hunter) RegisterSpells() {
	hunter.registerArcaneShotSpell()
	hunter.registerAspects()
	hunter.registerKillCommandSpell()
	hunter.registerMultiShotSpell()
	hunter.registerRaptorStrikeSpell()
	hunter.registerRapidFireCD()
	hunter.registerScorpidStingSpell()
	hunter.registerSerpentStingSpell()
	hunter.registerSteadyShotSpell()
	// hunter.registerHuntersMarkSpell()
}

func (hunter *Hunter) AddStatDependencies() {
	hunter.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.RangedAttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[hunter.Class])
	hunter.AddStatDependency(stats.Agility, stats.DodgeRating, 1.0/25*core.DodgeRatingPerDodgePercent)
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

func (hunter *Hunter) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if hunter.Talents.TrueshotAura {
		partyBuffs.TrueshotAura = true
	}
}

func (hunter *Hunter) Reset(_ *core.Simulation) {
}

func (hunter *Hunter) OnEncounterStart(sim *core.Simulation) {
}

const (
	HunterSpellFlagsNone int64 = 0
	SpellMaskSpellRanged int64 = 1 << iota
	HunterSpellAutoShot
	HunterSpellAimedShot
	HunterSpellArcaneShot
	HunterSpellAspectOfTheHawk
	HunterSpellAspectOfTheViper
	HunterSpellBestialWrath
	HunterSpellKillCommand
	HunterSpellKillCommandPet
	HunterSpellMultiShot
	HunterSpellRapidFire
	HunterSpellRaptorStrike
	HunterSpellRaptorStrikeQueue
	HunterSpellReadiness
	HunterSpellScorpidSting
	HunterSpellSerpentSting
	HunterSpellSteadyShot
	HunterSpellVolley
	HunterPetDamage
	HunterSpellsAll = HunterSpellAimedShot |
		HunterSpellArcaneShot | HunterSpellBestialWrath |
		HunterSpellKillCommand | HunterSpellMultiShot |
		HunterSpellRapidFire | HunterSpellRaptorStrike |
		HunterSpellScorpidSting | HunterSpellSerpentSting |
		HunterSpellSteadyShot | HunterSpellVolley
	HunterSpellsShotsAndStings = HunterSpellAimedShot |
		HunterSpellArcaneShot | HunterSpellMultiShot |
		HunterSpellScorpidSting | HunterSpellSerpentSting |
		HunterSpellSteadyShot | HunterSpellVolley
)

// Agent is a generic way to access underlying hunter on any of the agents.
type HunterAgent interface {
	GetHunter() *Hunter
}
