import {
	NodeService,
	SidebarItemKind,
	SidebarService,
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
} from '$lib/constants/constants';
import { SidebarAdd, SidebarLoad, SidebarRemove, SidebarUpdateName } from '$lib/constants/events';
import type { PanelType } from '$lib/state/layout-state.svelte';
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
		if (item.kind !== SidebarItemKind.SidebarItemKindGroup) {
			acc.push(item);
		}

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

		if (item.kind === SidebarItemKind.SidebarItemKindNetwork) {
			this.sidebar.root.name = item.name;

			const mappedItem = this.getItem(item.id);
			if (mappedItem) {
				mappedItem.name = item.name;
			}

			return;
		}

		const prefixes = item.prefix.split('/');

		let parent = this.sidebar.root;
		let remainingPrefix = prefixes.length;

		while (remainingPrefix > 0) {
			let foundParent = false;
			const currPrefix = prefixes[prefixes.length - remainingPrefix];

			for (const child of parent.children || []) {
				if (child.id !== currPrefix && child.id !== 'group/' + currPrefix) {
					continue;
				}

				parent = child;
				foundParent = true;

				break;
			}

			if (!foundParent) {
				return;
			}

			remainingPrefix = remainingPrefix - 1;
		}

		if (!parent.children) {
			return;
		}

		for (const child of parent.children) {
			if (child.id === item.id) {
				child.name = item.name;

				const mappedItem = this.getItem(item.id);
				if (mappedItem) {
					mappedItem.name = item.name;
				}

				break;
			}
		}

		parent.children.sort((a, b) => a.name.localeCompare(b.name));
	}

	private update(item: SidebarItem) {
		if (!this.sidebar) return;

		let parentItem = this.sidebar.root;
		for (const prefix of item.prefix.split('/')) {
			for (const child of parentItem.children || []) {
				const tmpSplitId = child.id.split('/');
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

	getKindOfGroup(groupId: string) {
		const splitId = groupId.split('/');
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
