<script lang="ts">
	import type { SignalUnit } from '$lib/api/canturin';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { SignalUnitIcon } from '$lib/components/icon';
	import type { PanelSectionProps } from '../types';
	import { getSignalUnitState } from './state.svelte';
	import { onMount } from 'svelte';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';

	let { entityId }: PanelSectionProps = $props();

	let sus = getSignalUnitState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await sus.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		nameSchema(() => invalidNames),
		() => sus.entity.name
	);

	function handleName(name: string) {
		sus.updateName(name);
	}

	function handleDesc(desc: string) {
		sus.updateDesc(desc);
	}
</script>

{#snippet section(sigUnit: SignalUnit)}
	<div class="flex gap-2 items-center">
		<SignalUnitIcon width="48" height="48" />

		<TextEditable
			bind:value={sigUnit.name}
			name="signal-unit-name"
			oncommit={handleName}
			errors={nameValidator.errors}
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={sigUnit.desc}
			name="signal-unit-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(sus.entity)}
</section>
