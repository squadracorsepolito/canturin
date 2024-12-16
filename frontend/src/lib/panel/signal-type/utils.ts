import { SignalTypeKind } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const signalTypeKindOptions: SegmentedControlOption[] = [
	{
		label: 'Custom',
		value: SignalTypeKind.SignalTypeKindCustom,
		desc: 'Any value'
	},
	{
		label: 'Flag',
		value: SignalTypeKind.SignalTypeKindFlag,
		desc: 'True or False'
	},
	{
		label: 'Integer',
		value: SignalTypeKind.SignalTypeKindInteger,
		desc: 'Any integer value'
	},
	{
		label: 'Decimal',
		value: SignalTypeKind.SignalTypeKindDecimal,
		desc: 'Any decimal value'
	}
];
