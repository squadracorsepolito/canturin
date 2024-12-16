<script lang="ts">
	import type { Node } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import { NumberEditable } from '$lib/components/editable';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';
	import { getNodeState } from './state.svelte';
	import { Readonly } from '$lib/components/readonly';
	import { getHexNumber } from '$lib/utils';
	import Divider from '$lib/components/divider/divider.svelte';

	let { entityId }: PanelSectionProps = $props();

	const ns = getNodeState(entityId);

	let invalidIds = $state<number[]>([]);

	async function loadInvalidIds() {
		const res = await ns.getInvalidIds();
		invalidIds = res;
	}

	$effect(() => {
		loadInvalidIds();
	});

	const idSchema = z.object({
		id: z
			.number()
			.min(0)
			.max(4_294_967_295)
			.refine((n) => !invalidIds.includes(n), { message: 'Duplicated' })
	});

	let idErrors = $derived.by(() => {
		const res = idSchema.safeParse({ id: ns.entity.id });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.id;
	});

	function handleId(id: number) {
		ns.updateID(id);
	}
</script>

{#snippet section(node: Node)}
	<div class="grid grid-cols-2 gap-5">
		<Attribute label="ID" desc="The ID of the node">
			<NumberEditable
				bind:value={node.id}
				name="node-id"
				errors={idErrors}
				oncommit={handleId}
				fontWeight="medium"
				border="visible"
			/>
		</Attribute>

		<Attribute label="Hex ID" desc="The ID of the node in hex">
			<Readonly>
				<span class="font-medium">{getHexNumber(node.id)}</span>
			</Readonly>
		</Attribute>
	</div>

	<Divider />

	<Attribute label="Interface Count" desc="The number of interfaces of the node">
		<Readonly>
			<span class="font-medium">{node.interfaces?.length}</span>
		</Readonly>
	</Attribute>
{/snippet}

<section>
	<h3 class="pb-5">Attributes</h3>

	{@render section(ns.entity)}
</section>
