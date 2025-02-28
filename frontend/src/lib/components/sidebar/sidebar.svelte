<script lang="ts">
	import { SidebarItemKind, type SidebarItem } from '$lib/api/canturin';
	import { deleteBus } from '$lib/panel/bus/state.svelte';
	import { deleteSignalEnum } from '$lib/panel/signal-enum/state.svelte';
	import { deleteSignalType } from '$lib/panel/signal-type/state.svelte';
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
		SignalUnitIcon
	} from '../icon';
	import { TreeView } from '../tree';
	import { SidebarState } from './state.svelte';
	import { deleteNode } from '$lib/panel/node/state.svelte';
	import { deleteSignalUnit } from '$lib/panel/signal-unit/state.svelte';
	import { AddSignalModal } from '../modal';
	import { IconButton } from '../button';

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

			default:
				return AltArrowIcon;
		}
	}

	function handleSelect(id: string) {
		s.openPanel(id);
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
			case SidebarItemKind.SidebarItemKindNodeInterface:
				s.addMessage();
				return;

			case SidebarItemKind.SidebarItemKindMessage:
				s.addMessage();
				return;

			case SidebarItemKind.SidebarItemKindSignalType:
			case SidebarItemKind.SidebarItemKindSignalUnit:
			case SidebarItemKind.SidebarItemKindSignalEnum:
		}
	}
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
			ondelete={handleDelete}
		>
			{#snippet actions({ collapse })}
				{#if s.selectedItemKind}
					{#if s.selectedItemKind === SidebarItemKind.SidebarItemKindSignal}
						<AddSignalModal onsubmit={(signalKind) => s.addSignal(signalKind)}>
							{#snippet trigger({ getProps })}
								<IconButton {...getProps()}>
									<AddIcon height="20" width="20" />
								</IconButton>
							{/snippet}
						</AddSignalModal>
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
</div>
