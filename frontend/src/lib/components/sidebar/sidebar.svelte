<script lang="ts">
	import { SidebarItemKind, type SidebarItem } from '$lib/api/canturin';
	import { deleteBus } from '$lib/panel/bus/state.svelte';
	import { deleteSignalEnum } from '$lib/panel/signal-enum/state.svelte';
	import layout from '$lib/state/layout-state.svelte';
	import { deleteSignalType } from '$lib/panel/signal-type/state.svelte';
	import {
		AltArrowIcon,
		BusIcon,
		MessageIcon,
		NetworkIcon,
		NodeIcon,
		SignalEnumIcon,
		SignalIcon,
		SignalTypeIcon,
		SignalUnitIcon
	} from '../icon';
	import { TreeView } from '../tree';
	import { SidebarState } from './state.svelte';
	import { deleteNode } from '$lib/panel/node/state.svelte';
	import { deleteSignalUnit } from '$lib/panel/signal-unit/state.svelte';

	const s = new SidebarState();

	function getTreeViewIcon(item: SidebarItem) {
		let kind = item.kind;

		if (kind === SidebarItemKind.SidebarItemKindGroup) {
			kind = s.getKindOfGroup(item.id);
		}

		switch (kind) {
			case SidebarItemKind.SidebarItemKindNetwork:
				return NetworkIcon;
			case SidebarItemKind.SidebarItemKindBus:
				return BusIcon;
			case SidebarItemKind.SidebarItemKindNode:
				return NodeIcon;
			case SidebarItemKind.SidebarItemKindMessage:
				return MessageIcon;
			case SidebarItemKind.SidebarItemKindSignal:
				return SignalIcon;
			case SidebarItemKind.SidebarItemKindSignalType:
				return SignalTypeIcon;
			case SidebarItemKind.SidebarItemKindSignalUnit:
				return SignalUnitIcon;
			case SidebarItemKind.SidebarItemKindSignalEnum:
				return SignalEnumIcon;

			default:
				return AltArrowIcon;
		}
	}

	function handleSelect(id: string) {
		const item = s.getItem(id);
		if (!item) return;

		layout.openPanel(s.getPanelType(item.kind), item.id);
	}

	function handleAdd(id: string) {
		const leafItem = s.getItem(id);
		if (leafItem) {
			layout.openPanel(s.getPanelType(leafItem.kind), 'draft');
			return;
		}

		// a group is selected
		const itemKind = s.getKindOfGroup(id);
		layout.openPanel(s.getPanelType(itemKind), 'draft');
	}

	function handleDelete(id: string) {
		const item = s.getItem(id);
		// it happens when a group is selected
		if (!item) return;

		switch (item.kind) {
			case SidebarItemKind.SidebarItemKindBus:
				deleteBus(item.id);
				break;

			case SidebarItemKind.SidebarItemKindNode:
				deleteNode(item.id);
				break;

			case SidebarItemKind.SidebarItemKindSignalType:
				deleteSignalType(item.id);
				break;

			case SidebarItemKind.SidebarItemKindSignalUnit:
				deleteSignalUnit(item.id);
				break;

			case SidebarItemKind.SidebarItemKindSignalEnum:
				deleteSignalEnum(item.id);
				break;
		}
	}

	$effect(() => {
		s.setSelectedItemId(layout.openPanelId);
	});
</script>

<div class="overflow-y-auto">
	{#if s.sidebar}
		<TreeView
			root={s.sidebar.root}
			bind:selectedValue={s.selectedItemId}
			valueKey="id"
			labelKey="name"
			getIcon={getTreeViewIcon}
			onselect={handleSelect}
			onadd={handleAdd}
			ondelete={handleDelete}
		/>
	{/if}
</div>
