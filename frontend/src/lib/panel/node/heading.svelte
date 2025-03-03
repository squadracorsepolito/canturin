<script lang="ts">
	import type { Node } from '$lib/api/canturin';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';
	import { getNodeState } from './state.svelte';
	import { NodeIcon } from '$lib/components/icon';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { onMount } from 'svelte';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';

	let { entityId }: PanelSectionProps = $props();

	let ns = getNodeState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await ns.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		nameSchema(() => invalidNames),
		() => ns.entity.name
	);

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
			errors={nameValidator.errors}
			oncommit={handleName}
			name="node-name"
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
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
