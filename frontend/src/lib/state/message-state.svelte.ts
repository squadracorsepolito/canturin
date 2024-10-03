import { MessageService, type Message } from '$lib/api/canturin';
import { AsyncState } from './async-state.svelte';

async function loadMessage(state: MessageState, entityId: string) {
	state.isLoading = true;

	try {
		const msg = await MessageService.Get(entityId);
		state.message = msg;
	} catch (error) {
		state.message = undefined;
		console.error(error);
	}

	state.isLoading = false;
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

const gridWidth = 8;

class MessageState extends AsyncState {
	message = $state<Message | undefined>();

	gridItems: GridItem[] = $derived.by(() => this.genGridItems());
	hoveredID = $state('');

	reload(entityId: string) {
		loadMessage(this, entityId);
	}

	private genGridItems() {
		if (!this.message || !this.message.signals) {
			return [];
		}

		const items: GridItem[] = [];

		for (const sig of this.message.signals) {
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
}

export function useMessage(entityId: string) {
	const state = new MessageState();
	loadMessage(state, entityId);
	return state;
}
