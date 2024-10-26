<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import * as numberInput from '@zag-js/number-input';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { untrack } from 'svelte';

	type Props = {
		initialValue: number;
		name: string;
		placeholder?: string;
		validator: (value: number) => string[] | undefined;
		onsubmit: (value: number) => void;
	};

	let { initialValue, name, placeholder, validator, onsubmit }: Props = $props();

	let errors = $state<string[]>();

	const inputId = uniqueId() + ':input';

	const [snpshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			value: initialValue + '',
			name: name,
			activationMode: 'dblclick',
			placeholder: placeholder,
			submitMode: 'both',
			autoResize: true,
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(initialValue + '');
					errors = undefined;
					return;
				}

				onsubmit(+details.value);
			},
			onValueChange: (details) => {
				errors = validator(parseFloat(details.value));
			},
			ids: {
				input: inputId
			}
		})
	);

	const [inputSnapshot, inputSend] = useMachine(
		numberInput.machine({
			id: uniqueId(),
			name: name,
			ids: {
				input: inputId
			}
		})
	);

	const api = $derived(editable.connect(snpshot, send, normalizeProps));
	const inputApi = $derived(numberInput.connect(inputSnapshot, inputSend, normalizeProps));

	$effect(() => {
		untrack(() => api.setValue)(initialValue + '');
	});
</script>

<div {...api.getRootProps()}>
	<div {...api.getAreaProps()} data-error={errors ? true : undefined}>
		<input {...inputApi.getInputProps()} {...api.getInputProps()} />

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

<style lang="postcss">
	[data-part='area'] {
		@apply rounded-btn px-3 py-2 border-2 inline-block font-medium transition-colors border-base-300;

		input {
			@apply outline-none;
		}

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
