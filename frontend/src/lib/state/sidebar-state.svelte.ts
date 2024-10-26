import { SidebarService } from '$lib/api/canturin';
import { SidebarAdd, SidebarLoad, SidebarRemove, SidebarUpdate } from '$lib/api/events';
import * as wails from '@wailsio/runtime';
import { SidebarNodeKind, type SidebarNode } from '$lib/api/canturin/models';
import type { WailsEvent } from '@wailsio/runtime/types/events';

class SidebarState {
	tree = $state<SidebarNode>();

	constructor() {
		wails.Events.On(SidebarLoad, () => {
			this.load();
		});

		wails.Events.On(SidebarUpdate, (e: WailsEvent) => {
			this.update(e.data[0] as SidebarNode);
		});

		wails.Events.On(SidebarAdd, (e: WailsEvent) => {
			console.log(e);

			this.update(e.data[0] as SidebarNode);
		});

		wails.Events.On(SidebarRemove, (e: WailsEvent) => {
			this.update(e.data[0] as SidebarNode);
		});

		// TODO! Remove this line in production
		this.load();
	}

	load() {
		const f = async () => {
			const rootNode = await SidebarService.GetTree();
			this.tree = rootNode;
		};

		f();
	}

	update(node: SidebarNode) {
		// console.log(node);

		let idx: number;
		switch (node.kind) {
			case SidebarNodeKind.SidebarNodeKindNetwork:
				this.tree = node;
				return;

			case SidebarNodeKind.SidebarNodeKindBus:
			case SidebarNodeKind.SidebarNodeKindSignalType:
			case SidebarNodeKind.SidebarNodeKindSignalUnit:
			case SidebarNodeKind.SidebarNodeKindSignalEnum:
				if (!this.tree || !this.tree.children) return;

				idx = this.tree.children.findIndex((n) => n.entityId === node.entityId);
				if (idx !== -1) this.tree.children[idx] = node;

				return;

			case SidebarNodeKind.SidebarNodeKindNode:
			case SidebarNodeKind.SidebarNodeKindMessage:
		}
	}
}

export default new SidebarState();
