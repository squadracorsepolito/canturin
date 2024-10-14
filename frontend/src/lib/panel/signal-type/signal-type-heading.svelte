<script lang="ts">
	import { SignalTypeIcon } from '$lib/components/icon';
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeState } from '../../state/signal-type-state.svelte';
	import { type SignalType } from '$lib/api/canturin';
	import { TextEditable } from '$lib/components/editable';
	import { z } from 'zod';
	import TextareaEditable from '$lib/components/editable/textarea-editable.svelte';
	import { text } from './signal-type-text';

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

	function handleName(name: string) {
		sts.updateName(name);
	}

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

	function validateName(name: string) {
		const res = nameSchema.safeParse({ name: name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	}

	function handleDesc(desc: string) {
		sts.updateDesc(desc);
	}
</script>

{#snippet section(signalType: SignalType)}
	<div class="flex gap-2 items-center">
		<SignalTypeIcon width="48" height="48" />

		<TextEditable
			validator={validateName}
			name="signal-type-name"
			initialValue={signalType.name}
			onsubmit={handleName}
			placeholder="Name"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={signalType.desc}
			name="signal-type-desc"
			triggerLabel={text.buttons.heading.descTriggerLabel}
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
