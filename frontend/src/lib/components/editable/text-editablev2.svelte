<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps } from '@zag-js/svelte';

	type Props = {
		value: string;
		name?: string;
		placeholder?: string;
		errors?: string[];
		fontSize?: 'lg' | 'md';
		oncommit?: (value: string) => void;
	};

	let {
		value = $bindable(),
		name,
		placeholder,
		errors,
		fontSize = 'md',
		oncommit
	}: Props = $props();

	let fallbackValue = $state(value);

	const [snapshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			name: name,
			activationMode: 'dblclick',
			placeholder: placeholder,
			autoResize: true,
			submitMode: 'both',
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(fallbackValue);
					return;
				}

				fallbackValue = details.value;
				oncommit?.(details.value);
			},
			onValueChange: (details) => {
				value = details.value;
			}
		}),
		{
			context: {
				get value() {
					return value;
				}
			}
		}
	);

	const api = $derived(editable.connect(snapshot, send, normalizeProps));
</script>

<div class="relative">
	<div {...api.getRootProps()}>
		<div {...api.getAreaProps()} data-error={errors ? true : undefined} data-font-size={fontSize}>
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

		input {
			@apply outline-none;
		}

		&[data-font-size='md'] {
			@apply text-h4;
		}

		&[data-font-size='lg'] {
			@apply text-h2;
		}
	}
</style>
