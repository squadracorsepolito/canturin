import { SignalTypeKind } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const signalTypeKindLables = {
	[SignalTypeKind.$zero]: '',
	[SignalTypeKind.SignalTypeKindCustom]: 'Custom',
	[SignalTypeKind.SignalTypeKindFlag]: 'Flag',
	[SignalTypeKind.SignalTypeKindInteger]: 'Integer',
	[SignalTypeKind.SignalTypeKindDecimal]: 'Decimal'
};

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
