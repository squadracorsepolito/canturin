import { SidebarService, type SidebarNode } from '$lib/api/canturin';
import { SidebarLoad, SidebarUpdate } from '$lib/api/events';
import * as wails from '@wailsio/runtime';

class SidebarState {
	tree = $state<SidebarNode>();

	constructor() {
		wails.Events.On(SidebarLoad, () => {
			this.load();
		});

		wails.Events.On(SidebarUpdate, () => {
			this.load();
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
}

export default new SidebarState();
