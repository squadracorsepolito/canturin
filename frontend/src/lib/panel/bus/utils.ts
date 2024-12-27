import { BusType } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export type BaudrateSelectItem = {
	label: string;
	value: string;
	valueAsNumber: number;
};

const baudrateValues = {
	'1_000_000': 1_000_000,
	'500_000': 500_000,
	'250_000': 250_000,
	'125_000': 125_000
};

function getSelectItemFromBaudrate(baudrate: number): BaudrateSelectItem {
	let bits = baudrateValues['1_000_000'];
	if (baudrate <= baudrateValues['125_000']) {
		bits = baudrateValues['125_000'];
	} else if (baudrate <= baudrateValues['250_000']) {
		bits = baudrateValues['250_000'];
	} else if (baudrate <= baudrateValues['500_000']) {
		bits = baudrateValues['500_000'];
	}

	const scale = bits < baudrateValues['1_000_000'] ? 1_000 : 1_000_000;
	const unit = scale === baudrateValues['1_000_000'] ? 'M' : 'k';

	return {
		label: `${bits / scale} ${unit}bit/s`,
		value: `${bits}`,
		valueAsNumber: bits
	};
}

export const baudrateItems = [
	getSelectItemFromBaudrate(125_000),
	getSelectItemFromBaudrate(250_000),
	getSelectItemFromBaudrate(500_000),
	getSelectItemFromBaudrate(1_000_000)
];

export const busTypeOptions: SegmentedControlOption[] = [
	{
		label: 'CAN 2.0 A',
		value: BusType.BusTypeCAN2A,
		desc: '11 bit identifier, max payload 8 bytes, max speed 1 Mbit/s'
	}
];
