import type { PanelType } from '$lib/state/layout-state.svelte';

type Panel = {
	kind: PanelType;
	id: string;
	name: string;
};

class PanelStackState {
	stack = $state<Panel[]>([]);
	diplayedIdx = $state(-1);

	displayedPanel = $derived.by(() => {
		if (this.diplayedIdx === -1) return undefined;

		return this.stack[this.diplayedIdx];
	});

	open(kind: PanelType, id: string, name: string) {
		this.stack.push({ kind, id, name });
		this.diplayedIdx = this.stack.length - 1;
	}

	close(panelIdx: number) {
		this.stack.splice(panelIdx, 1);

		if (panelIdx === this.diplayedIdx) {
			this.diplayedIdx = this.stack.length - 1;
		}
	}
}

const state = new PanelStackState();

export function getPanelStackState() {
	return state;
}

export function openPanel(kind: PanelType, id: string, name: string) {
	state.open(kind, id, name);
}

export function closePanel(panelIdx: number) {
	state.close(panelIdx);
}
