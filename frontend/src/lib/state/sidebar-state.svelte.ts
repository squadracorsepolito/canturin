import { SidebarService } from '$lib/api/canturin';
import { SidebarAdd, SidebarLoad, SidebarRemove, SidebarUpdate } from '$lib/api/events';
import { Events as wails } from '@wailsio/runtime';
import { SidebarNodeKind, type SidebarNode } from '$lib/api/canturin/models';

class SidebarState {
	tree = $state<SidebarNode>();

	constructor() {
		wails.On(SidebarLoad, () => {
			this.load();
		});

		wails.On(SidebarUpdate, (e: wails.WailsEvent) => {
			this.update(e.data[0] as SidebarNode);
		});

		wails.On(SidebarAdd, (e: wails.WailsEvent) => {
			console.log(e);

			this.update(e.data[0] as SidebarNode);
		});

		wails.On(SidebarRemove, (e: wails.WailsEvent) => {
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
