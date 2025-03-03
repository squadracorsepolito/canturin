<script lang="ts">
	import type { Bus } from '$lib/api/canturin';
	import { BusIcon } from '$lib/components/icon';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { onMount } from 'svelte';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';

	let { entityId }: PanelSectionProps = $props();

	let bs = getBusState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await bs.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		nameSchema(() => invalidNames),
		() => bs.entity.name
	);

	function handleName(name: string) {
		bs.updateName(name);
	}

	function handleDesc(desc: string) {
		bs.updateDesc(desc);
	}
</script>

{#snippet section(bus: Bus)}
	<div class="flex items-center gap-2">
		<BusIcon width="48" height="48" />

		<TextEditable
			bind:value={bus.name}
			errors={nameValidator.errors}
			oncommit={handleName}
			name="bus-name"
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={bus.desc}
			name="bus-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(bs.entity)}
</section>
