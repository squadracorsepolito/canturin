import { SignalTypeKind } from '$lib/api/canturin';

export const data = {
	name: {
		label: 'Name'
	},
	desc: {
		label: 'Description (optional)'
	},
	kind: {
		label: 'Kind',
		desc: 'The kind of the type',
		options: [
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
		]
	},
	size: {
		label: 'Size',
		desc: 'The size in bits'
	},
	signed: {
		label: 'Signed',
		desc: 'Whether the value is signed'
	},
	min: {
		label: 'Min',
		desc: 'The minimum value'
	},
	max: {
		label: 'Max',
		desc: 'The maximum value'
	},
	scale: {
		label: 'Scale',
		desc: 'The scale factor'
	},
	offset: {
		label: 'Offset',
		desc: 'The offset'
	},

	headings: {
		draft: 'Create new Signal Type'
	},

	buttons: {
		draft: {
			submit: 'Create Signal Type'
		}
	}
};
