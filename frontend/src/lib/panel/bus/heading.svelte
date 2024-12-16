<script lang="ts">
	import type { Bus } from '$lib/api/canturin';
	import { BusIcon } from '$lib/components/icon';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';
	import { getBusState } from './state.svelte';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';

	let { entityId }: PanelSectionProps = $props();

	let bs = getBusState(entityId);

	let invalidNames = $state<string[]>([]);

	async function loadInvalidNames() {
		const res = await bs.getInvalidNames();
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
		const res = nameSchema.safeParse({ name: bs.entity.name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	});

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
			errors={nameErrors}
			oncommit={handleName}
			name="bus-name"
			fontWeight="semibold"
			textSize="lg"
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
