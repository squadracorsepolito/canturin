<script lang="ts">
	import type { Network } from '$lib/api/canturin';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { NetworkIcon } from '$lib/components/icon';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';
	import { getNetworkState } from './state.svelte';

	let ns = getNetworkState();

	const nameValidator = new Validator(nameSchema(), () => ns.entity.name);

	function handleName(name: string) {
		ns.updateName(name);
	}

	function handleDesc(desc: string) {
		ns.updateDesc(desc);
	}
</script>

{#snippet section(net: Network)}
	<div class="flex gap-2 items-center">
		<NetworkIcon width="48" height="48" />

		<TextEditable
			bind:value={net.name}
			name="network-name"
			oncommit={handleName}
			errors={nameValidator.errors}
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={net.desc}
			name="network-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(ns.entity)}
</section>
