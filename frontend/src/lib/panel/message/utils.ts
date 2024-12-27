import { MessageSendType } from '$lib/api/canturin';

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
