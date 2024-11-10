<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps } from '@zag-js/svelte';

	type Props = {
		value: string;
		name?: string;
		placeholder?: string;
		errors?: string[];
		textSize?: 'md' | 'lg';
		fontWeight?: 'normal' | 'medium';
		oncommit?: (value: string) => void;
	};

	let {
		value = $bindable(),
		name,
		placeholder,
		errors,
		textSize = 'md',
		fontWeight = 'normal',
		oncommit
	}: Props = $props();

	let fallbackValue = $state(value);

	const [snapshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			name: name,
			activationMode: 'dblclick',
			placeholder: placeholder
				? {
						edit: '',
						preview: placeholder
					}
				: undefined,
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
		<div
			{...api.getAreaProps()}
			data-error={errors ? true : undefined}
			data-text-size={textSize}
			data-font-weight={fontWeight}
		>
			<input {...api.getInputProps()} />

			<span {...api.getPreviewProps()}>
				{api.valueText}
			</span>
		</div>
	</div>

	{#if errors && api.editing}
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
			@apply focus-ring-warning border-warning;

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
			@apply text-dimmed italic;
		}

		input {
			@apply outline-none bg-base-100;
		}

		&[data-text-size='lg'] {
			@apply text-h2;
		}

		&[data-font-weight='medium'] {
			@apply font-medium;
		}
	}

	[data-part='input'] {
		@apply outline-none bg-base-100;
	}
</style>
