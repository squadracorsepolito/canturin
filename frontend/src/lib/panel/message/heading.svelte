<script lang="ts">
	import { onMount } from 'svelte';
	import type { PanelSectionProps } from '../types';
	import { getMessageState } from './state.svelte';
	import * as v from 'valibot';
	import type { Message } from '$lib/api/canturin';
	import { MessageIcon } from '$lib/components/icon';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { checkUnused, Validator } from '$lib/utils/validator.svelte';
	import { Breadcrumbs } from '$lib/components/breadcrumbs';

	let { entityId }: PanelSectionProps = $props();

	const ms = getMessageState(entityId);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await ms.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		v.pipe(
			v.string(),
			v.minLength(1),
			checkUnused(() => invalidNames)
		),
		() => ms.entity.name
	);

	function handleName(name: string) {
		ms.updateName(name);
	}

	function handleDesc(desc: string) {
		ms.updateDesc(desc);
	}
</script>

{#snippet section(msg: Message)}
	{#if msg.paths}
		<Breadcrumbs paths={msg.paths} />
	{/if}

	<div class="flex gap-2 items-center pt-5">
		<MessageIcon width="48" height="48" />

		<TextEditable
			bind:value={msg.name}
			name="message-name"
			oncommit={handleName}
			errors={nameValidator.errors}
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
