<script lang="ts">
	import * as numberInput from '@zag-js/number-input';
	import { normalizeProps, useMachine } from '@zag-js/svelte';

	type Props = {
		value: number;
		name: string;
		label?: string;
		desc?: string;
		errors?: string[];
		min?: number;
		max?: number;
		disabled?: boolean;
	};

	let { value = $bindable(), name, label, desc, errors, min, max, disabled }: Props = $props();

	const id = $props.id();
	const service = useMachine(numberInput.machine, () => ({
		id,
		name,
		defaultValue: value.toString(),
		min,
		max,
		disabled,
		clampValueOnBlur: true,
		invalid: errors ? errors.length > 0 : false,
		onValueChange: (details) => {
			let tmpVal = details.valueAsNumber;
			if (isNaN(tmpVal)) {
				tmpVal = 0;
			}
			value = tmpVal;
		}
	}));

	const api = $derived(numberInput.connect(service, normalizeProps));
</script>

<div {...api.getRootProps()}>
	{#if label}
		<div class="flex-1">
			<label {...api.getLabelProps()}>
				<h4>{label}</h4>

				{#if desc}
					<p>{desc}</p>
				{/if}
			</label>
		</div>
	{/if}

	<div class="flex-1 relative">
		<div {...api.getControlProps()}>
			<input {...api.getInputProps()} />
		</div>

		{#if errors}
			<div data-scope="number-input" data-part="error">
				{#each errors as err}
					<span>{err}</span>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style lang="postcss">
	[data-scope='number-input'][data-part='root'] {
		@apply flex items-start gap-3;
	}

	[data-scope='number-input'][data-part='label'] {
		@apply flex-1;

		p {
			@apply text-sm text-dimmed;
		}

		&[data-invalid] {
			@apply text-error;
		}
	}

	[data-scope='number-input'][data-part='control'] {
		@apply rounded-btn border-2 px-2 py-1 border-neutral-content transition-colors;

		&[data-focus] {
			@apply border-primary focus-ring-primary bg-primary-ghost text-primary;

			&[data-invalid] {
				@apply focus-ring-error;
			}
		}

		&[data-invalid] {
			@apply border-error bg-error-ghost text-error;
		}
	}

	[data-scope='number-input'][data-part='input'] {
		@apply outline-none bg-transparent w-full;
	}

	[data-scope='number-input'][data-part='error'] {
		@apply absolute flex gap-1 pt-1;

		span {
			@apply truncate text-error text-xs;
		}
	}
</style>
