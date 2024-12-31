import { MessageByteOrder, MessageSendType, MessageService, type Message } from '$lib/api/canturin';
import { HistoryMessageModify } from '$lib/api/events';
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

export async function deleteMessage(entityId: string) {
	// TODO! implement
	console.log(entityId);
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

	deleteSignals(signalEntityIds: string[]) {
		this.update(MessageService.RemoveSignals(this.entity.entityId, { signalEntityIds }));
	}

	deleteSignal(signalEntityId: string) {
		this.deleteSignals([signalEntityId]);
	}

	compactSignals() {
		this.update(MessageService.CompactSignals(this.entity.entityId));
	}
}
