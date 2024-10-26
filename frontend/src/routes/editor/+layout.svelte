<script lang="ts">
	import { HistoryService, SidebarNodeKind, type SidebarNode } from '$lib/api/canturin';
	import { IconButton } from '$lib/components/button';
	import {
		NetworkIcon,
		BusIcon,
		NodeIcon,
		MessageIcon,
		SignalUnitIcon,
		SignalTypeIcon,
		AddIcon,
		SignalEnumIcon,
		UndoIcon,
		RedoIcon
	} from '$lib/components/icon';
	import Tree from '$lib/components/tree/tree.svelte';
	import type { TreeNode } from '$lib/components/tree/types';
	import layout from '$lib/state/layout-state.svelte';
	import sidebarState from '$lib/state/sidebar-state.svelte';
	import history from '$lib/state/history-state.svelte';

	let { children } = $props();

	function getNetTree(currNode: SidebarNode) {
		const n: TreeNode = {
			name: currNode.name,
			icon: NetworkIcon,
			childNodes: []
		};

		switch (currNode.kind) {
			case SidebarNodeKind.SidebarNodeKindNetwork:
				n.childNodes.push(
					{
						name: 'Signal Types',
						icon: SignalTypeIcon,
						childNodes: [
							{
								name: 'Add Signal Type',
								icon: AddIcon,
								childNodes: [],
								onclick: () => layout.openPanel('signal_type', 'draft')
							}
						]
					},
					{
						name: 'Signal Units',
						icon: SignalUnitIcon,
						childNodes: [
							{
								name: 'Add Signal Unit',
								icon: AddIcon,
								childNodes: [],
								onclick: () => layout.openPanel('signal_unit', 'draft')
							}
						]
					},
					{
						name: 'Signal Enums',
						icon: SignalEnumIcon,
						childNodes: [
							{
								name: 'Add Signal Enum',
								icon: AddIcon,
								childNodes: [],
								onclick: () => layout.openPanel('signal_enum', 'draft')
							}
						]
					}
				);
				break;

			case SidebarNodeKind.SidebarNodeKindBus:
				n.icon = BusIcon;
				break;

			case SidebarNodeKind.SidebarNodeKindNode:
				n.icon = NodeIcon;
				break;

			case SidebarNodeKind.SidebarNodeKindMessage:
				n.icon = MessageIcon;
				n.onclick = () => layout.openPanel('message', currNode.entityId);
				break;

			case SidebarNodeKind.SidebarNodeKindSignalType:
				n.icon = SignalTypeIcon;
				n.onclick = () => layout.openPanel('signal_type', currNode.entityId);
				break;

			case SidebarNodeKind.SidebarNodeKindSignalUnit:
				n.icon = SignalUnitIcon;
				n.onclick = () => layout.openPanel('signal_unit', currNode.entityId);
				break;

			case SidebarNodeKind.SidebarNodeKindSignalEnum:
				n.icon = SignalEnumIcon;
				n.onclick = () => layout.openPanel('signal_enum', currNode.entityId);
				break;
		}

		if (currNode.children) {
			for (const child of currNode.children) {
				const childNode = getNetTree(child);

				switch (child.kind) {
					case SidebarNodeKind.SidebarNodeKindSignalType:
						n.childNodes[0].childNodes.push(childNode);
						break;

					case SidebarNodeKind.SidebarNodeKindSignalUnit:
						n.childNodes[1].childNodes.push(childNode);
						break;

					case SidebarNodeKind.SidebarNodeKindSignalEnum:
						n.childNodes[2].childNodes.push(childNode);
						break;

					default:
						n.childNodes.push(childNode);
						break;
				}
			}
		}

		return n;
	}

	function handleUndo() {
		history.undo();
	}

	function handleRedo() {
		history.redo();
	}

	$inspect(history.history);
</script>

<div class="flex h-full w-full">
	<div class="h-full min-w-60 bg-base-200 flex flex-col">
		<div class="h-12 block bg-base-300 sticky top-0"></div>

		<div class="flex-1 overflow-y-auto overflow-x-hidden">
			{#if sidebarState.tree}
				<Tree rootNode={getNetTree(sidebarState.tree)} />
			{/if}
		</div>
	</div>

	<div class="flex-1 flex flex-col">
		<div class="h-12 bg-base-200 sticky top-0 block">
			<div class="flex items-center h-full px-5 gap-2">
				<IconButton onclick={handleUndo} disabled={!history.canUndo}>
					<UndoIcon />
				</IconButton>

				<IconButton onclick={handleRedo} disabled={!history.canRedo}>
					<RedoIcon />
				</IconButton>
			</div>
		</div>

		{@render children()}
	</div>
</div>
