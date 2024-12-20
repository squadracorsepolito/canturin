<script lang="ts">
	import { SignalTypeIcon } from '$lib/components/icon';
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeState } from './state.svelte';
	import { type SignalType } from '$lib/api/canturin';
	import { TextEditable } from '$lib/components/editable';
	import { z } from 'zod';
	import TextareaEditable from '$lib/components/editable/textarea-editable.svelte';

	let { entityId }: PanelSectionProps = $props();

	let sts = getSignalTypeState(entityId);

	let invalidNames = $state<string[]>([]);

	async function loadInvalidNames() {
		const res = await sts.getInvalidNames();
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
		const res = nameSchema.safeParse({ name: sts.entity.name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	});

	function handleName(name: string) {
		sts.updateName(name);
	}

	function handleDesc(desc: string) {
		sts.updateDesc(desc);
	}
</script>

{#snippet section(signalType: SignalType)}
	<div class="flex gap-2 items-center">
		<SignalTypeIcon width="48" height="48" />

		<TextEditable
			bind:value={signalType.name}
			name="signal-type-name"
			oncommit={handleName}
			errors={nameErrors}
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={signalType.desc}
			name="signal-type-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
