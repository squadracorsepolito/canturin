<script lang="ts">
	import { SidebarItemKind, type SidebarItem } from '$lib/api/canturin';
	import {
		AddIcon,
		AltArrowIcon,
		BusIcon,
		CollapseIcon,
		MessageIcon,
		NetworkIcon,
		NodeIcon,
		SignalEnumIcon,
		SignalIcon,
		SignalTypeIcon,
		SignalUnitIcon,
		CanIdBuilderIcon
	} from '../icon';
	import { TreeView } from '../tree';
	import { SidebarState } from './state.svelte';
	import { AddNodeModal, AddSignalModal, AddSignalTypeModal } from '../modal';
	import { IconButton } from '../button';
	import { openPanel } from '$lib/panel/panel-stack-state.svelte';

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
			case SidebarItemKind.SidebarItemKindNodeInterface:
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
			case SidebarItemKind.SidebarItemKindCANIDBuilder:
				return CanIdBuilderIcon;

			default:
				return AltArrowIcon;
		}
	}

	function handleSelect(id: string) {
		s.handleOpenPanel(id);
	}

	function handleDelete(id: string) {
		s.deleteItem(id);
	}

	function handleAdd() {
		if (!s.selectedItemKind) return;

		switch (s.selectedItemKind) {
			case SidebarItemKind.SidebarItemKindGroup:
			case SidebarItemKind.SidebarItemKindNetwork:
				return;

			case SidebarItemKind.SidebarItemKindBus:
				s.addBus();
				return;

			case SidebarItemKind.SidebarItemKindNode:
				return;

			case SidebarItemKind.SidebarItemKindNodeInterface:
			case SidebarItemKind.SidebarItemKindMessage:
				s.addMessage();
				return;

			case SidebarItemKind.SidebarItemKindSignalUnit:
				s.addSignalUnit();
				return;

			case SidebarItemKind.SidebarItemKindSignalEnum:
				s.addSignalEnum();
				return;
		}
	}

	$inspect(s.sidebar);
</script>

{#if s.sidebar}
	{@const netName = s.sidebar.root ? s.sidebar.root.name : ''}

	<TreeView
		root={s.sidebar.root}
		bind:selectedValue={s.selectedItemId}
		valueKey="id"
		labelKey="name"
		getIcon={getTreeViewIcon}
		onrootclick={() => openPanel('network', '', netName)}
		onselect={handleSelect}
		ondelete={handleDelete}
	>
		{#snippet actions({ collapse })}
			{#if s.selectedItemKind}
				{#if s.selectedItemKind === SidebarItemKind.SidebarItemKindNode}
					<AddNodeModal onsubmit={(nodeKind) => s.addNode(nodeKind)}>
						{#snippet trigger({ getProps })}
							<IconButton {...getProps()}>
								<AddIcon height="20" width="20" />
							</IconButton>
						{/snippet}
					</AddNodeModal>
				{:else if s.selectedItemKind === SidebarItemKind.SidebarItemKindSignal}
					<AddSignalModal onsubmit={(signalKind) => s.addSignal(signalKind)}>
						{#snippet trigger({ getProps })}
							<IconButton {...getProps()}>
								<AddIcon height="20" width="20" />
							</IconButton>
						{/snippet}
					</AddSignalModal>
				{:else if s.selectedItemKind === SidebarItemKind.SidebarItemKindSignalType}
					<AddSignalTypeModal
						onsubmit={(signalTypeKind, size) => s.addSignalType(signalTypeKind, size)}
					>
						{#snippet trigger({ getProps })}
							<IconButton {...getProps()}>
								<AddIcon height="20" width="20" />
							</IconButton>
						{/snippet}
					</AddSignalTypeModal>
				{:else}
					<IconButton onclick={handleAdd}>
						<AddIcon height="20" width="20" />
					</IconButton>
				{/if}
			{/if}

			<IconButton onclick={() => collapse()}>
				<CollapseIcon height="20" width="20" />
			</IconButton>
		{/snippet}
	</TreeView>
{/if}
