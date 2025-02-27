import {
	NodeService,
	SidebarItemKind,
	SidebarService,
	type Sidebar,
	type SidebarItem
} from '$lib/api/canturin';
import {
	SidebarNodeGroupID,
	SidebarSignalEnumGroupID,
	SidebarSignalTypeGroupID,
	SidebarSignalUnitGroupID
} from '$lib/constants/constants';
import { SidebarAdd, SidebarLoad, SidebarRemove, SidebarUpdateName } from '$lib/constants/events';
import type { PanelType } from '$lib/state/layout-state.svelte';
import layoutState from '$lib/state/layout-state.svelte';
import { Events as wails } from '@wailsio/runtime';

export class SidebarState {
	#items = $derived.by(() => {
		const m = new Map<string, SidebarItem>();

		if (this.sidebar) {
			const flattenedItems: SidebarItem[] = [];
			this.flattenItems(this.sidebar.root, flattenedItems);

			for (const item of flattenedItems) {
				m.set(item.id, item);
			}
		}

		return m;
	});

	sidebar = $state<Sidebar>();
	selectedItemId = $state('');

	constructor() {
		$effect(() => {
			this.setSelectedItemId(layoutState.openPanelId);
		});

		wails.On(SidebarLoad, () => {
			this.sidebar = undefined;
			this.load();
		});

		wails.On(SidebarUpdateName, (e: wails.WailsEvent) => {
			this.updateName(e.data[0] as SidebarItem);
		});

		wails.On(SidebarAdd, (e: wails.WailsEvent) => {
			this.update(e.data[0] as SidebarItem);
		});

		wails.On(SidebarRemove, (e: wails.WailsEvent) => {
			this.update(e.data[0] as SidebarItem);
		});

		// TODO!: remove this line in production
		this.load();
	}

	private flattenItems(item: SidebarItem, acc: SidebarItem[]) {
		acc.push(item);

		for (const child of item.children || []) {
			this.flattenItems(child, acc);
		}
	}

	async load() {
		const sidebar = await SidebarService.Get();
		this.sidebar = sidebar;
	}

	private updateName(item: SidebarItem) {
		if (!this.sidebar) return;

		const mappedItem = this.getItem(item.id);
		if (!mappedItem) return;

		mappedItem.name = item.name;

		const splPath = item.path.split('/');
		const splPtahLen = splPath.length;
		if (splPtahLen < 2) return;

		const parentId = splPath[splPtahLen - 2];
		const mappedParentItem = this.getItem(parentId);
		if (!mappedParentItem || !mappedParentItem.children) return;

		mappedParentItem.children.sort((a, b) => a.name.localeCompare(b.name));
	}

	private update(item: SidebarItem) {
		if (!this.sidebar) return;

		let parentItem = this.sidebar.root;

		const splPath = item.path.split('/');
		if (splPath[0] !== parentItem.id) {
			return;
		}

		for (let i = 1; i < splPath.length - 1; i++) {
			const pathId = splPath[i];

			for (const child of parentItem.children || []) {
				if (child.id !== pathId) continue;

				parentItem = child;
				break;
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

	getItem(id: string) {
		return this.#items.get(id);
	}

	getKindOfGroup(groupId: string) {
		switch (groupId) {
			case SidebarNodeGroupID:
				return SidebarItemKind.SidebarItemKindNode;
			case SidebarSignalTypeGroupID:
				return SidebarItemKind.SidebarItemKindSignalType;
			case SidebarSignalUnitGroupID:
				return SidebarItemKind.SidebarItemKindSignalUnit;
			case SidebarSignalEnumGroupID:
				return SidebarItemKind.SidebarItemKindSignalEnum;
		}

		const item = this.getItem(groupId);
		if (item) {
			return item.kind;
		}

		return SidebarItemKind.SidebarItemKindNetwork;
	}

	getPanelType(itemKind: SidebarItemKind): PanelType {
		switch (itemKind) {
			case SidebarItemKind.SidebarItemKindNetwork:
				break;
			case SidebarItemKind.SidebarItemKindBus:
				return 'bus';
			case SidebarItemKind.SidebarItemKindNode:
				return 'node';
			case SidebarItemKind.SidebarItemKindNodeInterface:
				return 'node';
			case SidebarItemKind.SidebarItemKindMessage:
				return 'message';
			case SidebarItemKind.SidebarItemKindSignal:
				return 'signal';
			case SidebarItemKind.SidebarItemKindSignalType:
				return 'signal_type';
			case SidebarItemKind.SidebarItemKindSignalUnit:
				return 'signal_unit';
			case SidebarItemKind.SidebarItemKindSignalEnum:
				return 'signal_enum';
		}

		return 'none';
	}

	setSelectedItemId(id: string) {
		this.selectedItemId = id;
	}

	openPanel(id: string) {
		const item = this.getItem(id);
		if (!item) return;

		if (item.kind === SidebarItemKind.SidebarItemKindGroup) return;

		const panelType = this.getPanelType(item.kind);

		let entityId = item.id;
		if (item.kind === SidebarItemKind.SidebarItemKindNodeInterface) {
			entityId = entityId.split(':')[0];
		}

		layoutState.openPanel(panelType, entityId);
	}

	async handleAdd(item: SidebarItem) {
		const splPath = item.path.split('/');
		const pathLen = splPath.length;

		switch (item.kind) {
			case SidebarItemKind.SidebarItemKindMessage:
				try {
					await NodeService.AddSentMessage(splPath[pathLen - 2], {
						interfaceNumber: +splPath[pathLen - 1]
					});
				} catch (err) {
					console.error(err);
				}

				break;
		}
	}
}
