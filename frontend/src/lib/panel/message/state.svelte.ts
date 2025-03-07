import {
	MessageByteOrder,
	MessageSendType,
	MessageService,
	NodeService,
	SignalKind,
	type Message
} from '$lib/api/canturin';
import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import { HistoryMessageModify } from '$lib/constants/events';
import { EntityState } from '$lib/state/entity-state.svelte';
import { StateProvider } from '$lib/state/state-provider.svelte';

const provider = new StateProvider(
	(message: Message) => new MessageState(message),
	HistoryMessageModify
);

export function getMessageState(entityId: string) {
	return provider.get(entityId);
}

export async function loadMessage(entityId: string) {
	const message = await MessageService.Get(entityId);
	provider.add(message);
}

export async function createMessage(nodeEntityId: string, interfaceNumber: number) {
	try {
		await NodeService.AddSentMessage(nodeEntityId, { interfaceNumber });
	} catch (err) {
		console.error(err);
		pushToast('error', 'Error', 'Operation failed');
	}
}

export async function deleteMessage(
	nodeEntityId: string,
	interfaceNumber: number,
	messageEntityId: string
) {
	try {
		await NodeService.RemoveSentMessages(nodeEntityId, {
			interfaceNumber,
			messageEntityIds: [messageEntityId]
		});
	} catch (err) {
		console.error(err);
		pushToast('error', 'Error', 'Operation failed');
	}
}

class MessageState extends EntityState<Message> {
	constructor(message: Message) {
		super(message);
	}

	async getInvalidNames() {
		const invalidNames = await MessageService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	async getInvalidMessageIds() {
		const invalidMsgIds = await MessageService.GetInvalidMessageIDs(this.entity.entityId);

		if (invalidMsgIds) {
			return invalidMsgIds;
		}

		return [];
	}

	async getInvalidCanIds() {
		const invalidCanIds = await MessageService.GetInvalidCANIDs(
			this.entity.entityId,
			this.entity.parentBus.entityId
		);

		if (invalidCanIds) {
			return invalidCanIds;
		}

		return [];
	}

	updateName(name: string) {
		this.update(MessageService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(MessageService.UpdateDesc(this.entity.entityId, { desc }));
	}

	updateMessageId(messageId: number) {
		this.update(MessageService.UpdateMessageID(this.entity.entityId, { messageId }));
	}

	updateStaticCanId(staticCanId: number) {
		this.update(MessageService.UpdateStaticCANID(this.entity.entityId, { staticCanId }));
	}

	updateSizeByte(sizeByte: number) {
		this.update(MessageService.UpdateSizeByte(this.entity.entityId, { sizeByte }));
	}

	updateByteOrder(byteOrder: MessageByteOrder) {
		this.update(MessageService.UpdateByteOrder(this.entity.entityId, { byteOrder }));
	}

	updateCycleTime(cycleTime: number) {
		this.update(MessageService.UpdateCycleTime(this.entity.entityId, { cycleTime }));
	}

	updateSendType(sendType: MessageSendType) {
		this.update(MessageService.UpdateSendType(this.entity.entityId, { sendType }));
	}

	updateDelayTime(delayTime: number) {
		this.update(MessageService.UpdateDelayTime(this.entity.entityId, { delayTime }));
	}

	updateStartDelayTime(startDelayTime: number) {
		this.update(MessageService.UpdateStartDelayTime(this.entity.entityId, { startDelayTime }));
	}

	addSignal(signalKind: SignalKind) {
		this.update(MessageService.AddSignal(this.entity.entityId, { signalKind }));
	}

	deleteSignals(signalEntityIds: string[]) {
		this.update(MessageService.DeleteSignals(this.entity.entityId, { signalEntityIds }));
	}

	deleteSignal(signalEntityId: string) {
		this.deleteSignals([signalEntityId]);
	}

	compactSignals() {
		this.update(MessageService.CompactSignals(this.entity.entityId));
	}

	reorderSignal(signalEntityId: string, from: number, to: number) {
		this.update(MessageService.ReorderSignal(this.entity.entityId, { signalEntityId, from, to }));
	}
}
