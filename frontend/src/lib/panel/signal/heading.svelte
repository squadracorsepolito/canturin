<script lang="ts">
	import type { Signal } from '$lib/api/canturin';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { SignalIcon } from '$lib/components/icon';
	import { onMount } from 'svelte';
	import type { PanelSectionProps } from '../types';
	import { getSignalState } from './state.svelte';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';

	let { entityId }: PanelSectionProps = $props();

	const ss = getSignalState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await ss.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		nameSchema(() => invalidNames),
		() => ss.entity.name
	);

	function handleName(name: string) {
		ss.updateName(name);
	}

	function handleDesc(desc: string) {
		ss.updateDesc(desc);
	}
</script>

{#snippet section(sig: Signal)}
	<div class="flex gap-2 items-center">
		<SignalIcon width="48" height="48" />

		<TextEditable
			bind:value={sig.name}
			errors={nameValidator.errors}
			oncommit={handleName}
			name="signal-name"
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={sig.desc}
			name="signal-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(ss.entity)}
</section>
