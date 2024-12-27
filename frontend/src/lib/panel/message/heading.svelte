<script lang="ts">
	import { onMount } from 'svelte';
	import type { PanelSectionProps } from '../types';
	import { getMessageState } from './state.svelte';
	import { z } from 'zod';
	import type { Message } from '$lib/api/canturin';
	import { MessageIcon } from '$lib/components/icon';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await ms.getInvalidNames();
		invalidNames = res;
	});

	const nameSchema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});

	let nameErrors = $derived.by(() => {
		const res = nameSchema.safeParse({ name: ms.entity.name });
		if (res.success) {
			return undefined;
		}
		return res.error.flatten().fieldErrors.name;
	});

	function handleName(name: string) {
		ms.updateName(name);
	}

	function handleDesc(desc: string) {
		ms.updateDesc(desc);
	}
</script>

{#snippet section(msg: Message)}
	<div class="flex gap-2 items-center">
		<MessageIcon width="48" height="48" />

		<TextEditable
			bind:value={msg.name}
			name="message-name"
			oncommit={handleName}
			errors={nameErrors}
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={msg.desc}
			name="message-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(ms.entity)}
</section>
