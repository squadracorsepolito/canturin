<script lang="ts">
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { z } from 'zod';
	import type { PanelSectionProps } from '../types';
	import { type SignalEnum } from '$lib/api/canturin';
	import { SignalEnumIcon } from '$lib/components/icon';
	import { TextEditable } from '$lib/components/editable';
	import { TextareaEditable } from '$lib/components/editable';
	import { text } from '../signal-type/signal-type-text';

	let { entityId }: PanelSectionProps = $props();

	let ses = getSignalEnumState(entityId);

	let invalidNames = $state<string[]>([]);

	async function loadInvalidNames() {
		const res = await ses.getInvalidNames();
		invalidNames = res;
	}

	$effect(() => {
		loadInvalidNames();
	});

	function handleName(name: string) {
		ses.updateName(name);
	}

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

	function validateName(name: string) {
		const res = nameSchema.safeParse(name);
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	}

	function handleDesc(desc: string) {
		ses.updateDesc(desc);
	}
</script>

{#snippet section(signalEnum: SignalEnum)}
	<div class="flex gap-2 items-center">
		<SignalEnumIcon width="48" height="48" />

		<TextEditable
			validator={validateName}
			name="signal-enum-name"
			initialValue={signalEnum.name}
			onsubmit={handleName}
			placeholder="Name"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={signalEnum.desc}
			name="signal-enum-desc"
			triggerLabel={text.buttons.heading.descTriggerLabel}
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(ses.entity)}
</section>
