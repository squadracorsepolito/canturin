<script lang="ts">
	import { NetworkService, type NetworkStub } from '$lib/api/canturin';
	import {
		NetworkIcon,
		BusIcon,
		NodeIcon,
		MessageIcon,
		SignalUnitIcon,
		SignalTypeIcon,
		AddIcon
	} from '$lib/components/icon';
	import Tree from '$lib/components/tree/tree.svelte';
	import type { TreeNode } from '$lib/components/tree/types';
	import { setLayoutState } from '$lib/state/layout-state.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();

	let net: NetworkStub | undefined = $state();

	onMount(async () => {
		try {
			const res = await NetworkService.GetNetworkStub();
			net = res;
		} catch (error) {
			console.error(error);
		}
	});

	const layout = setLayoutState();

	function getTreeNodes(network: NetworkStub) {
		let rootNode: TreeNode = {
			name: network.name,
			icon: NetworkIcon,
			childNodes: [],
			onclick: () => console.log('network')
		};

		if (network.signalTypes) {
			let sigTypes: TreeNode = {
				name: 'Signal Types',
				icon: SignalTypeIcon,
				childNodes: [],
				onclick: () => console.log('signal types')
			};

			sigTypes.childNodes.push({
				name: 'Add Signal Type',
				icon: AddIcon,
				childNodes: [],
				onclick: () => layout.openPanel('signal_type', 'draft')
			});

			for (let sigType of network.signalTypes) {
				let sigTypeNode: TreeNode = {
					name: sigType.name,
					icon: SignalTypeIcon,
					childNodes: [],
					onclick: () => layout.openPanel('signal_type', sigType.entityId)
				};
				sigTypes.childNodes.push(sigTypeNode);
			}

			rootNode.childNodes.push(sigTypes);
		}

		if (network.signalUnits) {
			let sigUnits: TreeNode = {
				name: 'Signal Units',
				icon: SignalUnitIcon,
				childNodes: [],
				onclick: () => console.log('signal units')
			};

			sigUnits.childNodes.push({
				name: 'Add Signal Unit',
				icon: AddIcon,
				childNodes: [],
				onclick: () => layout.openPanel('signal_unit', 'draft')
			});

			for (let sigUnit of network.signalUnits) {
				let sigUnitNode: TreeNode = {
					name: sigUnit.name,
					icon: SignalUnitIcon,
					childNodes: [],
					onclick: () => layout.openPanel('signal_unit', sigUnit.entityId)
				};
				sigUnits.childNodes.push(sigUnitNode);
			}

			rootNode.childNodes.push(sigUnits);
		}

		if (network.buses) {
			for (let bus of network.buses) {
				if (bus?.nodes) {
					let busNode: TreeNode = {
						name: bus.name,
						icon: BusIcon,
						childNodes: [],
						onclick: () => console.log('bus')
					};

					for (let node of bus.nodes) {
						if (node?.sendedMessages) {
							let nodeNode: TreeNode = {
								name: node.name,
								icon: NodeIcon,
								childNodes: [],
								onclick: () => console.log('node')
							};

							for (let sendMsg of node.sendedMessages) {
								if (sendMsg) {
									nodeNode.childNodes.push({
										name: sendMsg.name,
										icon: MessageIcon,
										childNodes: [],
										onclick: () => layout.openPanel('message', sendMsg.entityId)
									});
								}
							}

							busNode.childNodes.push(nodeNode);
						}
					}

					rootNode.childNodes.push(busNode);
				}
			}
		}

		return rootNode;
	}

	// async function openSignalUnit(sigEntId: string) {
	// 	try {
	// 		await SignalUnitService.Open(sigEntId);
	// 		entityId = sigEntId;
	// 		openPanel = 'signal_unit';
	// 	} catch (error) {
	// 		console.error(error);
	// 	}
	// }

	// async function openSignalType(sigTypeEntityId: string) {
	// 	try {
	// 		await SignalTypeService.Open(sigTypeEntityId);
	// 		entityId = sigTypeEntityId;
	// 		openPanel = 'signal_type';
	// 	} catch (error) {
	// 		console.error(error);
	// 	}
	// }

	// async function registerMessage(busEntID: string, nodeEntID: string, msgEntID: string) {
	// 	try {
	// 		await MessageService.Register(busEntID, nodeEntID, msgEntID);
	// 		const msg = await MessageService.Get();
	// 		message = msg;
	// 		entityId = msgEntID;

	// 		openPanel = 'message';
	// 	} catch (error) {
	// 		console.error(error);
	// 	}
	// }

	// let message: Message | undefined = $state();

	// let openPanel: 'signal_type' | 'signal_unit' | 'message' | 'none' = $state('none');
	// let entityId = $state('');
</script>

<div class="flex h-full w-full">
	<div class="h-full min-w-60 bg-base-200 flex flex-col">
		<div class="h-12 block bg-base-300 sticky top-0"></div>

		<div class="flex-1 overflow-y-auto overflow-x-hidden">
			{#if net}
				<Tree rootNode={getTreeNodes(net)} defaultOpen />
			{/if}
		</div>
	</div>

	<div class="flex-1 flex flex-col">
		<div class="h-12 block bg-base-300 sticky top-0"></div>

		{@render children()}

		<!-- {#if openPanel === 'message' && message}
			<MessagePanel {message} />
		{:else if openPanel === 'signal_unit'}
			<SignalUnitPanel {entityId} />
		{:else if openPanel === 'signal_type'}
			<SignalTypePanel {entityId} />
		{/if} -->
	</div>
</div>
