import { HistoryService, type History } from '$lib/api/canturin';
import { HistoryChange } from '$lib/api/events';
import * as wails from '@wailsio/runtime';
import type { WailsEvent } from '@wailsio/runtime/types/events';

class HistoryState {
	history = $state({
		operationCount: 0,
		currentIndex: -1
	});

	canUndo = $derived(this.history.operationCount > 0 && this.history.currentIndex > -1);
	canRedo = $derived(this.history.currentIndex < this.history.operationCount - 1);

	constructor() {
		wails.Events.On(HistoryChange, (event: WailsEvent) => {
			const history = event.data[0] as History;
			this.history = history;
		});
	}

	async undo() {
		const res = await HistoryService.Undo();
		this.history = res;
	}

	async redo() {
		const res = await HistoryService.Redo();
		this.history = res;
	}
}

export default new HistoryState();
