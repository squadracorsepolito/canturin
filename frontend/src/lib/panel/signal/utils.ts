import { SignalKind } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const signalKindOptions: SegmentedControlOption[] = [
	{
		label: 'Standard',
		value: SignalKind.SignalKindStandard,
		desc: 'Signal that has a type and can have a unit'
	},
	{
		label: 'Enum',
		value: SignalKind.SignalKindEnum,
		desc: 'Signal that represents an enum'
	},
	{
		label: 'Multiplexer',
		value: SignalKind.SignalKindMultiplexed,
		desc: 'Signal that contains multiple signals'
	}
];
