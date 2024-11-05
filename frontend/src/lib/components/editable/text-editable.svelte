<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { untrack } from 'svelte';

	type Props = {
		initialValue: string;
		name: string;
		placeholder?: string;
		size?: 'lg' | 'md';
		validator: (value: string) => string[] | undefined;
		onsubmit: (value: string) => void;
	};

	let { initialValue, name, placeholder, size = 'lg', validator, onsubmit }: Props = $props();

	let errors = $state<string[]>();

	const [snapshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			value: initialValue,
			name: name,
			activationMode: 'dblclick',
			placeholder: placeholder,
			autoResize: true,
			submitMode: 'both',
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(initialValue);
					errors = undefined;
					return;
				}

				onsubmit(details.value);
			},
			onValueChange: (details) => {
				errors = validator(details.value);
			}
		})
	);

	const api = $derived(editable.connect(snapshot, send, normalizeProps));

	$effect(() => {
		untrack(() => api.setValue)(initialValue);
	});

	function getFontSize() {
		switch (size) {
			case 'lg':
				return 'text-h2';
			case 'md':
			default:
				return 'text-h4';
		}
	}
</script>

<div>
	<div {...api.getRootProps()}>
		<div {...api.getAreaProps()} data-error={errors ? true : undefined} class={getFontSize()}>
			<input {...api.getInputProps()} />

			<span {...api.getPreviewProps()}>
				{api.valueText}
			</span>
		</div>
	</div>

	{#if errors}
		<div class="absolute pt-1 text-error text-xs">
			{#each errors as err}
				<span>{err}</span>
			{/each}
		</div>
	{/if}
</div>

<style lang="postcss">
	[data-part='area'] {
		@apply rounded-btn border-2 border-transparent px-2 py-1 transition-colors;

		&[data-error] {
			&[data-focus] {
				@apply focus-ring-error border-error;
			}
		}

		&:not([data-error]) {
			&[data-focus] {
				@apply focus-ring-primary border-primary;
			}
		}

		&[data-placeholder-shown] {
			@apply text-dimmed;
		}
	}
</style>
