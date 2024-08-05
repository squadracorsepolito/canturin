<script lang="ts">
	import { NetworkService, type Message, type Network } from '$lib/api/canturin';
	import { NetworkIcon, BusIcon, NodeIcon, MessageIcon } from '$lib/components/icon';
	import Tree from '$lib/components/tree/tree.svelte';
	import type { TreeNode } from '$lib/components/tree/types';
	import MessagePanel from '$lib/panel/message-panel.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();

	let net: Network | undefined = $state();

	onMount(async () => {
		const res = await NetworkService.GetNetwork();
		if (res) {
			net = res;
		}
	});

	function getTreeNodes(network: Network) {
		let rootNode: TreeNode = {
			name: network.name,
			icon: NetworkIcon,
			childNodes: [],
			onclick: () => console.log('network')
		};

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
										onclick: () => getMessage(bus.id, node.id, sendMsg.id)
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

	async function getMessage(busID: string, nodeID: string, msgID: string) {
		try {
			const res = await NetworkService.GetMessage(busID, nodeID, msgID);
			message = res;
		} catch (error) {
			console.error(error);
		}
	}

	let message: Message | undefined = $state();
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

		{#if message}
			<MessagePanel {message} />
		{/if}
	</div>
</div>
