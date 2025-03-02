import {
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
import { SidebarAdd, SidebarLoad, SidebarDelete, SidebarUpdateName } from '$lib/constants/events';
import type { PanelType } from '$lib/state/layout-state.svelte';
import layoutState from '$lib/state/layout-state.svelte';
import { Events as wails } from '@wailsio/runtime';
import { createBus, deleteBus } from '$lib/panel/bus/state.svelte';
import { createMessage, deleteMessage } from '$lib/panel/message/state.svelte';
import { createSignal, deleteSignal } from '$lib/panel/signal/state.svelte';

type SidebarUpdateNameEvent = {
	updatedId: string;
	name: string;
};

type SidebarAddEvent = {
	addedItem: SidebarItem;
};

type SidebarDeleteEvent = {
	deletedId: string;
};

export class SidebarState {
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

	#items = new Map<string, SidebarItem>();

	constructor() {
		$effect(() => {
			this.setSelectedItemId(layoutState.openPanelId);
		});

		wails.On(SidebarLoad, () => {
			this.sidebar = undefined;
			this.load();
		});

		wails.On(SidebarUpdateName, (e: wails.WailsEvent) => {
			this.handleUpdateNameEvent(e.data[0] as SidebarUpdateNameEvent);
		});

		wails.On(SidebarAdd, (e: wails.WailsEvent) => {
			this.handleAddEvent(e.data[0] as SidebarAddEvent);
		});

		wails.On(SidebarDelete, (e: wails.WailsEvent) => {
			this.handleDeleteEvent(e.data[0] as SidebarDeleteEvent);
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

	private getItem(id: string) {
		return this.#items.get(id);
	}

	private sortChildren(chlidren: SidebarItem[]) {
		chlidren.sort((a, b) => {
			if (a.kind === SidebarItemKind.SidebarItemKindGroup) {
				if (b.kind === SidebarItemKind.SidebarItemKindGroup) return 0;

				return -1;
			}

			if (b.kind === SidebarItemKind.SidebarItemKindGroup) return 1;

			return a.name.localeCompare(b.name);
		});
	}

	async load() {
		const sidebar = await SidebarService.Get();
		this.sidebar = sidebar;

		const flattenedItems: SidebarItem[] = [];
		this.flattenItems(this.sidebar.root, flattenedItems);

		for (const item of flattenedItems) {
			this.#items.set(item.id, item);
		}
	}

	private handleUpdateNameEvent(e: SidebarUpdateNameEvent) {
		const item = this.getItem(e.updatedId);
		if (!item) return;

		item.name = e.name;

		const parentId = this.getParentId(item.path);
		if (!parentId) return;

		const parent = this.getItem(parentId);
		if (!parent || !parent.children) return;

		const tmpChildren = [...parent.children];
		this.sortChildren(tmpChildren);
		parent.children = tmpChildren;
	}

	private handleAddEvent(e: SidebarAddEvent) {
		const item = e.addedItem;

		this.#items.set(item.id, item);
		for (const child of item.children || []) {
			this.#items.set(child.id, child);
		}

		const parentId = this.getParentId(item.path);
		if (!parentId) return;

		const parent = this.getItem(parentId);
		if (!parent) return;

		const newChildren: SidebarItem[] = [item, ...(parent.children || [])];
		this.sortChildren(newChildren);
		parent.children = newChildren;
	}

	private handleDeleteEvent(e: SidebarDeleteEvent) {
		const item = this.getItem(e.deletedId);
		if (!item) return;

		for (const child of item.children || []) {
			this.#items.delete(child.id);
		}

		this.#items.delete(item.id);

		const parentId = this.getParentId(item.path);
		if (!parentId) return;

		const parent = this.getItem(parentId);
		if (!parent || !parent.children) return;

		const newChildren: SidebarItem[] = [];
		for (const child of parent.children) {
			if (child.id !== item.id) {
				newChildren.push(child);
			}
		}
		parent.children = newChildren;
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

	private getParentId(path: string) {
		const splPath = path.split('/');
		if (splPath.length < 2) return '';

		return splPath[splPath.length - 2];
	}

	private getMessageParent(parentId: string) {
		const splParentId = parentId.split(':');

		return {
			nodeEntityId: splParentId[0],
			interfaceNumber: parseInt(splParentId[1])
		};
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

		let parentId = item.id;
		if (item.kind === SidebarItemKind.SidebarItemKindMessage) {
			parentId = this.getParentId(item.path);
		}

		const msgParent = this.getMessageParent(parentId);
		await createMessage(msgParent.nodeEntityId, msgParent.interfaceNumber);
	}

	async addSignal(signalKind: SignalKind) {
		const item = this.getItem(this.selectedItemId);
		if (!item) return;

		if (item.kind !== SidebarItemKind.SidebarItemKindSignal) return;

		await createSignal(this.getParentId(item.path), signalKind);
	}

	async addBus() {
		const item = this.getItem(this.selectedItemId);
		if (!item) return;

		if (item.kind !== SidebarItemKind.SidebarItemKindBus) return;

		await createBus();
	}

	async deleteItem(id: string) {
		const item = this.getItem(id);
		if (!item) return;

		switch (item.kind) {
			case SidebarItemKind.SidebarItemKindGroup:
			case SidebarItemKind.SidebarItemKindNetwork:
				return;

			case SidebarItemKind.SidebarItemKindBus:
				await deleteBus(item.id);
				return;

			case SidebarItemKind.SidebarItemKindNode:
				return;

			case SidebarItemKind.SidebarItemKindNodeInterface:
				return;

			case SidebarItemKind.SidebarItemKindMessage: {
				const msgParent = this.getMessageParent(this.getParentId(item.path));
				await deleteMessage(msgParent.nodeEntityId, msgParent.interfaceNumber, item.id);
				return;
			}

			case SidebarItemKind.SidebarItemKindSignal:
				await deleteSignal(this.getParentId(item.path), item.id);
				return;

			case SidebarItemKind.SidebarItemKindSignalType:
			case SidebarItemKind.SidebarItemKindSignalUnit:
			case SidebarItemKind.SidebarItemKindSignalEnum:
		}
	}
}
