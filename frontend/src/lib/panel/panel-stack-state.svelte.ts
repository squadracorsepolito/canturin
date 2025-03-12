import { EntityKind } from '$lib/api/canturin';
import { SvelteMap } from 'svelte/reactivity';

export type PanelKind =
	| 'none'
	| 'network'
	| 'bus'
	| 'node'
	| 'message'
	| 'signal'
	| 'signal_type'
	| 'signal_unit'
	| 'signal_enum';

export type Panel = {
	kind: PanelKind;
	id: string;
	name: string;
	prevId: string;
};

class PanelStackState {
	panels = new SvelteMap<string, Panel>();
	displayedPanel = $state<Panel>();

	open(kind: PanelKind, id: string, name: string) {
		let prevId = '';
		if (this.displayedPanel) {
			prevId = this.displayedPanel.id;
		}

		if (prevId === id) return;

		const panel: Panel = { kind, id, name, prevId };

		if (!this.panels.has(panel.id)) {
			this.panels.set(id, panel);
		}

		this.displayedPanel = panel;
	}

	close(panelId: string) {
		if (!this.panels.delete(panelId)) return;

		const dispPanel = this.displayedPanel;
		if (!dispPanel) return;

		if (dispPanel.id !== panelId) return;

		this.displayedPanel = this.panels.get(dispPanel.prevId);
	}
}

const state = new PanelStackState();

export function getPanelStackState() {
	return state;
}

export function openPanel(kind: PanelKind, id: string, name: string) {
	state.open(kind, id, name);
}

export function closePanel(panelId: string) {
	state.close(panelId);
}

export function getPanelKind(entityKind: EntityKind) {
	let kind: PanelKind;
	switch (entityKind) {
		case EntityKind.EntityKindNetwork:
			kind = 'network';
			break;
		case EntityKind.EntityKindBus:
			kind = 'bus';
			break;
		case EntityKind.EntityKindNode:
			kind = 'node';
			break;
		case EntityKind.EntityKindMessage:
			kind = 'message';
			break;
		case EntityKind.EntityKindSignal:
			kind = 'signal';
			break;
		case EntityKind.EntityKindSignalType:
			kind = 'signal_type';
			break;
		case EntityKind.EntityKindSignalUnit:
			kind = 'signal_unit';
			break;
		case EntityKind.EntityKindSignalEnum:
			kind = 'signal_enum';
			break;
		default:
			kind = 'network';
	}
	return kind;
}
