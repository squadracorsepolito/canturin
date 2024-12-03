<script lang="ts">
	import { SidebarItemKind, type SidebarItem } from '$lib/api/canturin';
	import layout from '$lib/state/layout-state.svelte';
	import {
		AltArrowIcon,
		BusIcon,
		MessageIcon,
		NetworkIcon,
		NodeIcon,
		SignalEnumIcon,
		SignalTypeIcon,
		SignalUnitIcon
	} from '../icon';
	import TreeView from '../tree/tree-view.svelte';
	import { SidebarState } from './state.svelte';

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

	let selectedId = $state('');

	function handleSelect(id: string) {
		const item = s.getItem(id);
		if (!item) return;

		layout.openPanel(s.getPanelType(item.kind), item.id);
	}

	function handleAdd() {
		if (!selectedId) return;

		const leafItem = s.getItem(selectedId);
		if (leafItem) {
			layout.openPanel(s.getPanelType(leafItem.kind), 'draft');
			return;
		}

		// a group is selected
		const itemKind = s.getKindOfGroup(selectedId);
		layout.openPanel(s.getPanelType(itemKind), 'draft');
	}
</script>

<div class="overflow-y-auto">
	{#if s.sidebar}
		<TreeView
			root={s.sidebar.root}
			bind:selectedValue={selectedId}
			valueKey="id"
			labelKey="name"
			getIcon={getTreeViewIcon}
			onselect={handleSelect}
			onadd={handleAdd}
		/>
	{/if}
</div>
