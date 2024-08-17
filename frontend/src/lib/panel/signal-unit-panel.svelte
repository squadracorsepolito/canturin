<script lang="ts">
	import { SignalUnitService, type SignalUnit, type SignalReference } from '$lib/api/canturin';
	import { ReferenceTree } from '$lib/components/tree';
	import type { ReferenceTreeNode } from '$lib/components/tree/types';
	import Panel from './panel.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let sigUnitPromise = $derived.by(async () => {
		try {
			const tmpSigUnit = await SignalUnitService.GetOpen(entityId);
			return tmpSigUnit;
		} catch (error) {
			console.error(error);
		}
	});

	function getReferences(refs: SignalReference[]): ReferenceTreeNode[] {
		const res: ReferenceTreeNode[] = [];

		for (let ref of refs) {
			const signalNode: ReferenceTreeNode = {
				name: ref.signal.name,
				childNodes: []
			};

			const msgNode: ReferenceTreeNode = {
				name: ref.message.name,
				childNodes: [signalNode]
			};

			const nodeNode: ReferenceTreeNode = {
				name: ref.node.name,
				childNodes: [msgNode]
			};

			const busNode: ReferenceTreeNode = {
				name: ref.bus.name,
				childNodes: [nodeNode]
			};

			const bus = res.find((b) => b.name === ref.bus.name);
			if (bus === undefined) {
				res.push(busNode);
				continue;
			}

			const node = bus.childNodes.find((n) => n.name === ref.node.name);
			if (node === undefined) {
				bus.childNodes.push(nodeNode);
				continue;
			}

			const msg = node.childNodes.find((m) => m.name === ref.message.name);
			if (msg === undefined) {
				node.childNodes.push(msgNode);
				continue;
			}

			msg.childNodes.push(signalNode);
		}

		return res;
	}
</script>

{#snippet sigUnitPanel(sigUnit: SignalUnit)}
	<section>
		<h3>{sigUnit.name}</h3>
		<p>{sigUnit.desc}</p>
	</section>

	{#if sigUnit.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getReferences(sigUnit.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#await sigUnitPromise then sigUnit}
		{#if sigUnit}
			{@render sigUnitPanel(sigUnit)}
		{/if}
	{/await}
</Panel>
