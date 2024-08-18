<script lang="ts">
	import {
		MessageService,
		NetworkService,
		SignalTypeService,
		SignalUnitService,
		type Message,
		type NetworkStub
	} from '$lib/api/canturin';
	import {
		NetworkIcon,
		BusIcon,
		NodeIcon,
		MessageIcon,
		SignalUnitIcon,
		SignalTypeIcon
	} from '$lib/components/icon';
	import Tree from '$lib/components/tree/tree.svelte';
	import type { TreeNode } from '$lib/components/tree/types';
	import MessagePanel from '$lib/panel/message-panel.svelte';
	import SignalTypePanel from '$lib/panel/signal-type-panel.svelte';
	import SignalUnitPanel from '$lib/panel/signal-unit-panel.svelte';
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

			for (let sigtype of network.signalTypes) {
				let sigTypeNode: TreeNode = {
					name: sigtype.name,
					icon: SignalTypeIcon,
					childNodes: [],
					onclick: () => openSignalType(sigtype.entityId)
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

			for (let sigUnit of network.signalUnits) {
				let sigUnitNode: TreeNode = {
					name: sigUnit.name,
					icon: SignalUnitIcon,
					childNodes: [],
					onclick: () => openSignalUnit(sigUnit.entityId)
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
										onclick: () => {
											registerMessage(bus.entityId, node.entityId, sendMsg.entityId);
										}
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

	async function openSignalUnit(sigEntId: string) {
		try {
			await SignalUnitService.Open(sigEntId);
			entityId = sigEntId;
			openPanel = 'signal_unit';
		} catch (error) {
			console.error(error);
		}
	}

	async function openSignalType(sigTypeEntityId: string) {
		try {
			await SignalTypeService.Open(sigTypeEntityId);
			entityId = sigTypeEntityId;
			openPanel = 'signal_type';
		} catch (error) {
			console.error(error);
		}
	}

	async function registerMessage(busEntID: string, nodeEntID: string, msgEntID: string) {
		try {
			await MessageService.Register(busEntID, nodeEntID, msgEntID);
			const msg = await MessageService.Get();
			message = msg;
			entityId = msgEntID;

			openPanel = 'message';
		} catch (error) {
			console.error(error);
		}
	}

	let message: Message | undefined = $state();

	let openPanel: 'signal_type' | 'signal_unit' | 'message' | 'none' = $state('none');
	let entityId = $state('');
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
		<!-- {@render children()} -->
		<div class="h-12 block bg-base-300 sticky top-0"></div>

		{#if openPanel === 'message' && message}
			<MessagePanel {message} />
		{:else if openPanel === 'signal_unit'}
			<SignalUnitPanel {entityId} />
		{:else if openPanel === 'signal_type'}
			<SignalTypePanel {entityId} />
		{/if}
	</div>
</div>
