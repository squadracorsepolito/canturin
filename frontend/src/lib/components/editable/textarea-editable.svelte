<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { AddIcon } from '../icon';
	import { untrack } from 'svelte';

	type Props = {
		initialValue: string;
		name: string;
		triggerLabel: string;
		rows?: number;
		onsubmit: (value: string) => void;
	};

	let { initialValue, name, triggerLabel, rows = 8, onsubmit }: Props = $props();

	const service = useMachine(editable.machine, {
		id: uniqueId(),
		value: initialValue,
		name: name,
		activationMode: 'dblclick',
		submitMode: 'both',
		onInteractOutside: (e) => {
			console.log(e);
		},
		onValueCommit: (details) => {
			onsubmit(details.value);
		}
	});

	const api = $derived(editable.connect(service, normalizeProps));

	$effect(() => {
		untrack(() => api.setValue)(initialValue);
	});
</script>

<div {...api.getRootProps()}>
	<div {...api.getAreaProps()}>
		<textarea {...api.getInputProps() as any} {rows}></textarea>

		<span {...api.getPreviewProps()}>
			{api.value}
		</span>
	</div>

	{#if api.empty && !api.editing}
		<button {...api.getEditTriggerProps()} class="btn btn-sm btn-outline btn-primary border-2">
			<AddIcon />
			{triggerLabel}
		</button>
	{/if}
</div>

<style lang="postcss">
	[data-part='area'] {
		@apply border-base-300 rounded-btn transition-colors;

		textarea {
			@apply w-full resize-none outline-none bg-base-100;
		}

		&:not([data-placeholder-shown]) {
			@apply border-2 px-3 py-2;
		}

		&[data-focus] {
			@apply border-2 px-3 py-2 border-primary focus-ring-primary;
		}
	}
</style>
