import {
	MessageService,
	NodeService,
	SidebarItemKind,
	SidebarService,
	SignalKind,
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
import { pushToast } from '../toast/toast-provider.svelte';

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
	selectedItemKind = $derived.by(() => {
		const item = this.getItem(this.selectedItemId);
		if (!item) return;

		if (item.kind === SidebarItemKind.SidebarItemKindGroup) {
			return this.getKindOfGroup(item.id);
		}

		return item.kind;
	});

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

		if (
			item.kind === SidebarItemKind.SidebarItemKindGroup ||
			item.kind === SidebarItemKind.SidebarItemKindNodeInterface
		) {
			return;
		}

		const panelType = this.getPanelType(item.kind);
		layoutState.openPanel(panelType, item.id);
	}

	async addMessage() {
		const item = this.getItem(this.selectedItemId);
		if (!item) return;

		if (
			item.kind !== SidebarItemKind.SidebarItemKindMessage &&
			item.kind !== SidebarItemKind.SidebarItemKindNodeInterface
		) {
			return;
		}

		let nodeIntItemId = item.id;

		if (item.kind === SidebarItemKind.SidebarItemKindMessage) {
			const splPath = item.path.split('/');
			nodeIntItemId = splPath[splPath.length - 2];
		}

		const splParentId = nodeIntItemId.split(':');
		const nodeEntId = splParentId[0];
		const interfaceNumber = parseInt(splParentId[1]);

		try {
			await NodeService.AddSentMessage(nodeEntId, { interfaceNumber });
		} catch (err) {
			console.error(err);
			pushToast('error', 'Error', 'Operation failed');
		}
	}

	async addSignal(signalKind: SignalKind) {
		const item = this.getItem(this.selectedItemId);
		if (!item) return;

		if (item.kind !== SidebarItemKind.SidebarItemKindSignal) return;

		const splPath = item.path.split('/');
		const messageEntId = splPath[splPath.length - 2];

		try {
			await MessageService.AddSignal(messageEntId, { signalKind });
		} catch (err) {
			console.error(err);
			pushToast('error', 'Error', 'Operation failed');
		}
	}

	async addBus() {
		const item = this.getItem(this.selectedItemId);
		if (!item) return;

		if (item.kind !== SidebarItemKind.SidebarItemKindBus) return;
	}
}
