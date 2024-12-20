import { SignalUnitKind } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const signalUnitKindOptions: SegmentedControlOption[] = [
	{
		label: 'Custom',
		value: SignalUnitKind.SignalUnitKindCustom,
		desc: 'Not restricted to a specific domain'
	},
	{
		label: 'Temperature',
		value: SignalUnitKind.SignalUnitKindTemperature,
		desc: 'The unit represents a temperature quantity'
	},
	{
		label: 'Electrical',
		value: SignalUnitKind.SignalUnitKindElectrical,
		desc: 'The unit represents an electrical quantity'
	},
	{
		label: 'Power',
		value: SignalUnitKind.SignalUnitKindPower,
		desc: 'The unit represents a power quantity'
	}
];
