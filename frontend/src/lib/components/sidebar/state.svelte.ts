import {
	SidebarItemKind,
	SidebarService0,
	type Sidebar,
	type SidebarItem
} from '$lib/api/canturin';
import {
	SidebarBusesPrefix,
	SidebarMessagesPrefix,
	SidebarNodesPrefix,
	SidebarSignalEnumsPrefix,
	SidebarSignalTypesPrefix,
	SidebarSignalUnitsPrefix
} from '$lib/api/constants';
import { SidebarAdd, SidebarLoad, SidebarUpdateName } from '$lib/api/events';
import type { PanelType } from '$lib/state/layout-state.svelte';
import { Events as wails } from '@wailsio/runtime';

export class SidebarState {
	#items = new Map<string, SidebarItem>();

	sidebar = $state<Sidebar>();

	constructor() {
		wails.On(SidebarLoad, () => {
			this.load();
		});

		wails.On(SidebarUpdateName, (e: wails.WailsEvent) => {
			this.updateName(e.data[0] as SidebarItem);
		});

		wails.On(SidebarAdd, (e: wails.WailsEvent) => {
			this.add(e.data[0] as SidebarItem);
		});

		// TODO! Remove this line in production
		this.load();
	}

	private loadItem(item: SidebarItem) {
		if (item.kind !== SidebarItemKind.SidebarItemKindGroup) {
			this.#items.set(item.id, item);
		}

		if (item.children) {
			for (const child of item.children) {
				this.loadItem(child);
			}
		}
	}

	async load() {
		const sidebar = await SidebarService0.Get();
		this.sidebar = sidebar;

		this.loadItem(sidebar.root);
	}

	add(item: SidebarItem) {
		if (!this.sidebar) return;

		let parentItem = this.sidebar.root;
		for (const prefix of item.prefix.split(':')) {
			for (const child of parentItem.children || []) {
				const tmpSplitId = child.id.split(':');
				if (tmpSplitId.length < 2) {
					continue;
				}

				if (tmpSplitId[1] === prefix) {
					parentItem = child;
					break;
				}
			}
		}

		if (!parentItem.children) return;

		for (let i = 0; i < parentItem.children.length; i++) {
			if (parentItem.children[i].id === item.id) {
				parentItem.children[i] = item;
				break;
			}
		}
	}

	private updateName(item: SidebarItem) {
		if (!this.sidebar) return;

		if (item.kind === SidebarItemKind.SidebarItemKindNetwork) {
			this.sidebar.root.name = item.name;

			const mappedItem = this.getItem(item.id);
			if (mappedItem) {
				mappedItem.name = item.name;
			}

			return;
		}

		let parentItem = this.sidebar.root;
		for (const prefix of item.prefix.split(':')) {
			for (const child of parentItem.children || []) {
				const tmpSplitId = child.id.split(':');
				if (tmpSplitId.length < 2) {
					continue;
				}

				if (tmpSplitId[1] === prefix) {
					parentItem = child;
					break;
				}
			}
		}

		for (const child of parentItem.children || []) {
			if (child.id === item.id) {
				child.name = item.name;

				const mappedItem = this.getItem(item.id);
				if (mappedItem) {
					mappedItem.name = item.name;
				}

				break;
			}
		}
	}

	getKindOfGroup(groupId: string) {
		const splitId = groupId.split(':');
		const id = splitId[splitId.length - 1];

		switch (id) {
			case SidebarBusesPrefix:
				return SidebarItemKind.SidebarItemKindBus;
			case SidebarNodesPrefix:
				return SidebarItemKind.SidebarItemKindNode;
			case SidebarMessagesPrefix:
				return SidebarItemKind.SidebarItemKindMessage;
			case SidebarSignalTypesPrefix:
				return SidebarItemKind.SidebarItemKindSignalType;
			case SidebarSignalUnitsPrefix:
				return SidebarItemKind.SidebarItemKindSignalUnit;
			case SidebarSignalEnumsPrefix:
				return SidebarItemKind.SidebarItemKindSignalEnum;
		}

		const item = this.getItem(id);
		if (item) {
			return item.kind;
		}

		return SidebarItemKind.SidebarItemKindNetwork;
	}

	getItem(id: string) {
		return this.#items.get(id);
	}

	getPanelType(itemKind: SidebarItemKind): PanelType {
		switch (itemKind) {
			case SidebarItemKind.SidebarItemKindNetwork:
				break;
			case SidebarItemKind.SidebarItemKindBus:
				return 'bus';
			case SidebarItemKind.SidebarItemKindNode:
				return 'node';
			case SidebarItemKind.SidebarItemKindMessage:
				return 'message';
			case SidebarItemKind.SidebarItemKindSignalType:
				return 'signal_type';
			case SidebarItemKind.SidebarItemKindSignalUnit:
				return 'signal_unit';
			case SidebarItemKind.SidebarItemKindSignalEnum:
				return 'signal_enum';
		}

		return 'none';
	}
}
