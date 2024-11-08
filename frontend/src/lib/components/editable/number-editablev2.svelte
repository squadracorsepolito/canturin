<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import * as numberInput from '@zag-js/number-input';
	import { useMachine, normalizeProps, mergeProps } from '@zag-js/svelte';

	type Props = {
		value: number;
		name?: string;
		placeholder?: string;
		errors?: string[];
		oncommit?: (value: number) => void;
	};

	const inputId = uniqueId() + ':input';

	let { value = $bindable(), name, placeholder, errors, oncommit }: Props = $props();

	let fallbackValue = $state(value);

	const [snpshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			name: name,
			activationMode: 'dblclick',
			placeholder: placeholder,
			submitMode: 'both',
			autoResize: true,
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(fallbackValue + '');
					value = fallbackValue;
					return;
				}

				const numValue = +details.value;

				fallbackValue = numValue;
				oncommit?.(numValue);
			},
			onValueRevert(details) {
				value = +details.value;
			},
			ids: {
				input: inputId
			}
		}),
		{
			context: {
				get value() {
					return value + '';
				}
			}
		}
	);

	const [inputSnapshot, inputSend] = useMachine(
		numberInput.machine({
			id: uniqueId(),
			name: name,
			onValueChange(details) {
				if (isNaN(details.valueAsNumber)) {
					value = 0;
					return;
				}

				value = details.valueAsNumber;
			},
			ids: {
				input: inputId
			}
		}),
		{
			context: {
				get value() {
					return value + '';
				}
			}
		}
	);

	const api = $derived(editable.connect(snpshot, send, normalizeProps));
	const inputApi = $derived(numberInput.connect(inputSnapshot, inputSend, normalizeProps));

	const inputProps = $derived(mergeProps(inputApi.getInputProps(), api.getInputProps()));
</script>

<div class="relative">
	<div {...api.getRootProps()}>
		<div {...api.getAreaProps()} data-error={errors ? true : undefined}>
			<input {...inputProps} />

			<span {...api.getPreviewProps()}>
				{api.valueText}
			</span>
		</div>
	</div>

	{#if errors}
		<div class="absolute pt-1 text-error text-xs truncate">
			{#each errors as err}
				<span>{err}</span>
			{/each}
		</div>
	{/if}
</div>

<style lang="postcss">
	[data-part='area'] {
		@apply rounded-btn px-2 py-1 border-2 border-transparent transition-colors;

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
