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

export type GridItem = {
	id: string;
	colStart: number;
	colEnd: number;
	rowStart: number;
	rowEnd: number;
	name: string;
	continues?: boolean;
	follows?: boolean;
};

export const gridWidth = 8;

class MessageState extends EntityState<Message> {
	gridItems: GridItem[] = $derived.by(() => this.genGridItems());
	hoveredID = $state('');

	constructor(message: Message) {
		super(message);
	}

	private genGridItems() {
		if (!this.entity || !this.entity.signals) {
			return [];
		}

		const items: GridItem[] = [];

		for (const sig of this.entity.signals) {
			const startPos = sig.startPos;
			const size = sig.size;
			const name = sig.name;
			const id = sig.entityId;

			const nRows = Math.ceil(((startPos % gridWidth) + size) / gridWidth);

			const startRow =
				startPos % gridWidth === 0
					? Math.ceil(startPos / gridWidth) + 1
					: Math.ceil(startPos / gridWidth);

			const startCol = (startPos % gridWidth) + 1;

			if (startPos % 8 === 0 && size % 8 === 0) {
				// spans exactly over one/multiple rows
				items.push({
					colStart: 1,
					colEnd: gridWidth + 1,
					rowStart: startRow,
					rowEnd: startRow + nRows,
					name,
					id
				});
				continue;
			}

			if (nRows === 1) {
				items.push({
					colStart: startCol,
					colEnd: startCol + size,
					rowStart: startRow,
					rowEnd: startRow,
					name,
					id
				});
				continue;
			}

			if (nRows === 2) {
				items.push(
					{
						colStart: startCol,
						colEnd: gridWidth + 1,
						rowStart: startRow,
						rowEnd: startRow,
						name,
						id,
						continues: true
					},
					{
						colStart: 1,
						colEnd: size - (gridWidth - (startCol - 1)) + 1,
						rowStart: startRow + nRows - 1,
						rowEnd: startRow + nRows - 1,
						name,
						id,
						follows: true
					}
				);
				continue;
			}

			items.push({
				colStart: startCol,
				colEnd: gridWidth + 1,
				rowStart: startRow,
				rowEnd: startRow,
				name,
				id,
				continues: true
			});

			items.push({
				colStart: 1,
				colEnd: gridWidth + 1,
				rowStart: startRow + 1,
				rowEnd: startRow + nRows - 1,
				name,
				id,
				continues: true,
				follows: true
			});

			items.push({
				colStart: 1,
				colEnd: size - (startCol - 1) - gridWidth * (nRows - 2) + 1,
				rowStart: startRow + nRows - 1,
				rowEnd: startRow + nRows - 1,
				name,
				id,
				follows: true
			});
		}

		// flip colStart/colEnd because 0 is on the right
		for (let i = 0; i < items.length; i++) {
			const cs = items[i].colStart;
			if (cs < 5) {
				items[i].colStart = cs + 2 * (5 - cs);
			} else if (cs) {
				items[i].colStart = cs - 2 * (cs - 5);
			}

			const ce = items[i].colEnd;
			if (cs < 5) {
				items[i].colEnd = ce + 2 * (5 - ce);
			} else if (cs) {
				items[i].colEnd = ce - 2 * (ce - 5);
			}
		}

		return items;
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
}
