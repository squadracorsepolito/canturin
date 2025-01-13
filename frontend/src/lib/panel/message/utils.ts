import { MessageByteOrder, MessageSendType, SignalKind } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const sendTypeSelectItems = [
	{
		label: 'Unset',
		value: MessageSendType.MessageSendTypeUnset
	},
	{
		label: 'Cyclic',
		value: MessageSendType.MessageSendTypeCyclic
	},
	{
		label: 'Cyclic if Active',
		value: MessageSendType.MessageSendTypeCyclicIfActive
	},
	{
		label: 'Cyclic and Triggered',
		value: MessageSendType.MessageSendTypeCyclicAndTriggered
	},
	{
		label: 'Cyclic if Active and Triggered',
		value: MessageSendType.MessageSendTypeCyclicIfActiveAndTriggered
	}
];

export const byteOrderOptions: SegmentedControlOption[] = [
	{
		value: MessageByteOrder.MessageByteOrderLittleEndian,
		label: 'Little Endian',
		desc: 'Least significant byte first'
	},
	{
		value: MessageByteOrder.MessageByteOrderBigEndian,
		label: 'Big Endian',
		desc: 'Most significant byte first'
	}
];

export function getSignalKindString(kind: SignalKind) {
	switch (kind) {
		case SignalKind.SignalKindStandard:
			return 'Standard';
		case SignalKind.SignalKindEnum:
			return 'Enum';
		case SignalKind.SignalKindMultiplexed:
			return 'Multiplexed';
		default:
			return '';
	}
}
