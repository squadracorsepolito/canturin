<script lang="ts">
	import * as radioGroup from '@zag-js/radio-group';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { TickIcon } from '../icon';
	import { uniqueId } from '$lib/utils';
	import type { SegmentedControlOption } from './types';

	type Props = {
		selectedValue: string;
		name: string;
		options: SegmentedControlOption[];
		readOnly?: boolean;
		onchange?: (value: string) => void;
	};

	let { selectedValue = $bindable(), name, options, readOnly, onchange }: Props = $props();

	const radioGroupProps: radioGroup.Props = $derived({
		id: uniqueId(),
		name,
		orientation: 'horizontal',
		value: selectedValue,
		readOnly,
		onValueChange: (details) => {
			if (details.value) {
				selectedValue = details.value;
				onchange?.(details.value);
			}
		}
	});

	const service = useMachine(radioGroup.machine, () => radioGroupProps);

	const api = $derived(radioGroup.connect(service, normalizeProps));
</script>

<div {...api.getRootProps()} style:grid-template-columns="repeat({options.length}, minmax(0, 1fr))">
	{#each options as opt}
		<label {...api.getItemProps({ value: opt.value })}>
			<div {...api.getItemTextProps({ value: opt.value })}>
				<span>{opt.label}</span>

				{#if opt.value === api.value}
					<span>
						<TickIcon height={20} width={20} />
					</span>
				{/if}
			</div>

			{#if opt.desc}
				<span class="text-sm text-dimmed">{opt.desc}</span>
			{/if}

			<input {...api.getItemHiddenInputProps({ value: opt.value })} />
		</label>
	{/each}
</div>

<style lang="postcss">
	[data-part='root'] {
		@apply grid gap-x-2;
	}

	[data-part='item'] {
		@apply grid row-span-2 grid-rows-subgrid border-2 py-3 rounded-btn gap-1 transition-colors px-3;

		&[data-readonly] {
			&:not([data-state='checked']) {
				@apply bg-base-300 border-base-300 opacity-80;
			}
		}

		&:not([data-readonly]) {
			@apply border-base-300 cursor-pointer;

			&:not([data-state='checked']) {
				&:hover {
					@apply border-secondary text-secondary bg-secondary-ghost;
				}
			}
		}

		&[data-state='checked'] {
			@apply border-primary text-primary bg-primary-ghost;
		}

		&[data-focus] {
			@apply focus-ring-primary;
		}
	}

	[data-part='item-text'] {
		@apply font-medium flex items-center gap-1 justify-between overflow-x-hidden;

		span:first-child {
			@apply truncate;
		}
	}
</style>
