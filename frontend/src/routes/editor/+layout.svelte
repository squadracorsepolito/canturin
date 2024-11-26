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
	import { PaneGroup, Pane, PaneResizer } from 'paneforge';
	import TickIcon from '$lib/components/icon/tick-icon.svelte';

	let { children } = $props();

	function getNetTree(currNode: SidebarNode) {
		const indexes = {
			// bus: 0,
			// node: 1,
			signalType: 0,
			signalUnit: 1,
			signalEnum: 2
		};

		const n: TreeNode = {
			name: currNode.name,
			icon: NetworkIcon,
			childNodes: []
		};

		switch (currNode.kind) {
			case SidebarNodeKind.SidebarNodeKindNetwork:
				// n.childNodes[indexes.bus] = {
				// 	name: 'Buses',
				// 	icon: BusIcon,
				// 	childNodes: [
				// 		{
				// 			name: 'Add Bus',
				// 			icon: AddIcon,
				// 			childNodes: [],
				// 			onclick: () => layout.openPanel('bus', 'draft')
				// 		}
				// 	]
				// };

				// n.childNodes[indexes.node] = {
				// 	name: 'Nodes',
				// 	icon: NodeIcon,
				// 	childNodes: [
				// 		{
				// 			name: 'Add Node',
				// 			icon: AddIcon,
				// 			childNodes: [],
				// 			onclick: () => layout.openPanel('node', 'draft')
				// 		}
				// 	]
				// };

				// n.childNodes[indexes.signalType] = {
				// 	name: 'Signal Types',
				// 	icon: SignalTypeIcon,
				// 	childNodes: [
				// 		{
				// 			name: 'Add Signal Type',
				// 			icon: AddIcon,
				// 			childNodes: [],
				// 			onclick: () => layout.openPanel('signal_type', 'draft')
				// 		}
				// 	]
				// };

				// n.childNodes[indexes.signalUnit] = {
				// 	name: 'Signal Units',
				// 	icon: SignalUnitIcon,
				// 	childNodes: [
				// 		{
				// 			name: 'Add Signal Unit',
				// 			icon: AddIcon,
				// 			childNodes: [],
				// 			onclick: () => layout.openPanel('signal_unit', 'draft')
				// 		}
				// 	]
				// };

				// n.childNodes[indexes.signalEnum] = {
				// 	name: 'Signal Enums',
				// 	icon: SignalEnumIcon,
				// 	childNodes: [
				// 		{
				// 			name: 'Add Signal Enum',
				// 			icon: AddIcon,
				// 			childNodes: [],
				// 			onclick: () => layout.openPanel('signal_enum', 'draft')
				// 		}
				// 	]
				// };

				// break;

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
				n.childNodes.push({
					name: 'Open Bus',
					icon: TickIcon,
					childNodes: [],
					onclick: () => layout.openPanel('bus', currNode.entityId)
				});
				break;

			case SidebarNodeKind.SidebarNodeKindNode:
				n.icon = NodeIcon;
				n.childNodes.push({
					name: 'Open Node',
					icon: TickIcon,
					childNodes: [],
					onclick: () => layout.openPanel('node', currNode.entityId)
				});
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
					// case SidebarNodeKind.SidebarNodeKindBus:
					// 	n.childNodes[indexes.bus].childNodes.push(childNode);
					// 	break;

					// case SidebarNodeKind.SidebarNodeKindNode:
					// 	n.childNodes[indexes.node].childNodes.push(childNode);
					// 	break;

					case SidebarNodeKind.SidebarNodeKindSignalType:
						n.childNodes[indexes.signalType].childNodes.push(childNode);
						break;

					case SidebarNodeKind.SidebarNodeKindSignalUnit:
						n.childNodes[indexes.signalUnit].childNodes.push(childNode);
						break;

					case SidebarNodeKind.SidebarNodeKindSignalEnum:
						n.childNodes[indexes.signalEnum].childNodes.push(childNode);
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

	$inspect(sidebarState.tree);
</script>

<PaneGroup direction="horizontal" class="h-full w-full">
	<Pane defaultSize={15} class="h-full bg-base-200 flex flex-col">
		<div class="h-12 block bg-base-300 sticky top-0"></div>

		<div class="flex-1 overflow-y-auto overflow-x-hidden">
			{#if sidebarState.tree}
				<Tree rootNode={getNetTree(sidebarState.tree)} />
			{/if}
		</div>
	</Pane>

	<PaneResizer
		class="h-full w-1 bg-base-200 data-[active=pointer]:bg-accent hover:bg-accent transition-colors delay-75"
	></PaneResizer>

	<Pane class="flex-1 flex flex-col">
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
	</Pane>
</PaneGroup>
