import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, PseudoStat, Stat } from '../../core/proto/common';
import { Warlock_Options as WarlockOptions, WarlockOptions_Armor, WarlockOptions_CurseOptions, WarlockOptions_Summon } from '../../core/proto/warlock';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import BlankAPL from './apls/blank.apl.json'
import BlankGear from './gear_sets/blank.gear.json';
import PreRaid from './gear_sets/preraid.json';
import PreRaidFire from './gear_sets/destro_fire_preraid.json';
import T4Set from './gear_sets/t4.json';
import T4Fire from './gear_sets/destro_fire_t4.json';
import T5Set from './gear_sets/t5.json';
import T6Set from './gear_sets/t6.json';
import ZASet from './gear_sets/za.json';
import SWPSet from './gear_sets/swp.json';
import AfflictionRot from './apls/affliction.json';
import DemoRot from './apls/demonology.json';
import DestroRot from './apls/destruction.json';
import DestroFireRot from './apls/destro_fire.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BLANK_APL = PresetUtils.makePresetAPLRotation('Blank', BlankAPL)

export const BLANK_GEARSET = PresetUtils.makePresetGear('Blank', BlankGear);

export const PRE_RAID = PresetUtils.makePresetGear('Pre-Raid', PreRaid)
export const PRE_RAID_FIRE = PresetUtils.makePresetGear('Pre-Raid Fire', PreRaidFire)

export const T4 = PresetUtils.makePresetGear('T4', T4Set)
export const T4_FIRE = PresetUtils.makePresetGear('T4 - Fire', T4Fire)

export const T5 = PresetUtils.makePresetGear('T5', T5Set)
export const T6 = PresetUtils.makePresetGear('T6', T6Set)
export const ZA = PresetUtils.makePresetGear("Zul'Aman", ZASet)
export const SWP = PresetUtils.makePresetGear("Sunwell Plateau", SWPSet)


// Preset options for EP weights
export const DEFAULT_EP = PresetUtils.makePresetEpWeights(
	'Sub',
	Stats.fromMap(
		{
			[Stat.StatIntellect]: 0.4,
			[Stat.StatSpellDamage]: 1,
			[Stat.StatFireDamage]: 1,
			[Stat.StatShadowDamage]: 1,
			[Stat.StatSpellHitRating]: 0,
			[Stat.StatSpellCritRating]: 0.8,
			[Stat.StatSpellHasteRating]: 1.2,
		}
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const Talents = {
	name: 'A',
	data: SavedTalents.create({
		talentsString: '',
	}),
};

export const Affliction = {
	name: 'Affliction',
	data: SavedTalents.create({
		talentsString: '05022221112351055003--50500051220001',
	}),
};

export const DemoRuin = {
	name: 'Demo/Ruin',
	data: SavedTalents.create({
		talentsString: '01-205003213305010150134-50500251020001',
	})
}

export const DemoFelguard = {
	name: 'Demonology Felguard',
	data: SavedTalents.create({
		talentsString: '01-2050030133250101501351-5050005112'
	})
}

export const DestroNightfall = {
	name: 'Destro/Nightfall',
	data: SavedTalents.create({
		talentsString: '150222201023--505020510200510531051',
	})
}

export const Destruction = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '-20500301332101-50500051220051053105'
	})
}

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		armor: WarlockOptions_Armor.FelArmor,
		curseOptions: WarlockOptions_CurseOptions.Agony,
		sacrificeSummon: true,
		summon: WarlockOptions_Summon.Succubus,
	},
});

// Rotations
export const AfflictionAPL = PresetUtils.makePresetAPLRotation("Affliction", AfflictionRot)
export const DemoAPL = PresetUtils.makePresetAPLRotation("Demonology", DemoRot)
export const DestroAPL = PresetUtils.makePresetAPLRotation("Destruction", DestroRot)
export const DestroFireAPL = PresetUtils.makePresetAPLRotation("Destro - Fire", DestroFireRot)

// Defaults

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22866, // Flask of Pure Death
	foodId: 27657, // Blackened Basilisk
	mhImbueId: 20749, // Brilliant Wizard Oil
	petScrollAgi: true,
	petScrollStr: true,
	potId: 22839, // Destruction Potion
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
