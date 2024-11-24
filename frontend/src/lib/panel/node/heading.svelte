<script lang="ts">
	import type { Node } from '$lib/api/canturin';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';
	import { getNodeState } from './state.svelte';
	import { NodeIcon } from '$lib/components/icon';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';

	let { entityId }: PanelSectionProps = $props();

	let ns = getNodeState(entityId);

	let invalidNames = $state<string[]>([]);

	async function loadInvalidNames() {
		const res = await ns.getInvalidNames();
		invalidNames = res;
	}

	$effect(() => {
		loadInvalidNames();
	});

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

	let nameErrors = $derived.by(() => {
		const res = nameSchema.safeParse({ name: ns.entity.name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	});

	function handleName(name: string) {
		ns.updateName(name);
	}

	function handleDesc(desc: string) {
		ns.updateDesc(desc);
	}
</script>

{#snippet section(node: Node)}
	<div class="flex items-center gap-2">
		<NodeIcon width="48" height="48" />

		<TextEditable
			bind:value={node.name}
			errors={nameErrors}
			oncommit={handleName}
			name="node-name"
			fontWeight="semibold"
			textSize="lg"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={node.desc}
			name="node-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(ns.entity)}
</section>
