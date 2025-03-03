<script lang="ts">
	import { getSignalEnumState } from '$lib/panel/signal-enum/state.svelte';
	import type { PanelSectionProps } from '../types';
	import { type SignalEnum } from '$lib/api/canturin';
	import { SignalEnumIcon } from '$lib/components/icon';
	import { TextEditable } from '$lib/components/editable';
	import { TextareaEditable } from '$lib/components/editable';
	import { onMount } from 'svelte';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';

	let { entityId }: PanelSectionProps = $props();

	let ses = getSignalEnumState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await ses.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		nameSchema(() => invalidNames),
		() => ses.entity.name
	);

	function handleName(name: string) {
		ses.updateName(name);
	}

	function handleDesc(desc: string) {
		ses.updateDesc(desc);
	}
</script>

{#snippet section(signalEnum: SignalEnum)}
	<div class="flex gap-2 items-center">
		<SignalEnumIcon width="48" height="48" />

		<TextEditable
			bind:value={signalEnum.name}
			name="signal-enum-name"
			oncommit={handleName}
			errors={nameValidator.errors}
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={signalEnum.desc}
			name="signal-enum-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(ses.entity)}
</section>
