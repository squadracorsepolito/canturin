import { MessageByteOrder, MessagePriority, MessageSendType } from '$lib/api/canturin';
import type { SegmentedControlOption } from '$lib/components/segmented-control/types';

export const prioritySelectItems = [
	{
		label: 'Very High',
		value: MessagePriority.MessagePriorityVeryHigh
	},
	{
		label: 'High',
		value: MessagePriority.MessagePriorityHigh
	},
	{
		label: 'Medium',
		value: MessagePriority.MessagePriorityMedium
	},
	{
		label: 'Low',
		value: MessagePriority.MessagePriorityLow
	}
];

export function priorityToNumber(priority: MessagePriority) {
	switch (priority) {
		case MessagePriority.MessagePriorityVeryHigh:
			return 0;
		case MessagePriority.MessagePriorityHigh:
			return 1;
		case MessagePriority.MessagePriorityMedium:
			return 2;
		case MessagePriority.MessagePriorityLow:
			return 3;
		default:
			return 0;
	}
}

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
